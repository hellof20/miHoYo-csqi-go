package main

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"cloud.google.com/go/vertexai/genai"
	"github.com/tidwall/gjson"
)

//go:embed label_rules.json
var labelRules []byte

type Record struct {
	ID      string
	Result  string
	Comment string
	Label   string
	Content string
}

type FirstCatalog struct {
	Label  string `json:"label"`
	Reason string `json:"reason"`
}

type SecondCatalog struct {
	PrimaryClassification   string `json:"primary_classification"`
	SecondaryClassification string `json:"secondary_classification"`
	Reason                  string `json:"reason"`
}

type Result struct {
	ID           string `json:"id"`
	First_label  string `json:"first_label"`
	Second_label string `json:"second_label"`
	Thrid_label  string `json:"thrid_label"`
	Reason       string `json:"reason"`
}

type LLMAPI interface {
	InvokeText(prompt string) (string, error)
	SetResponseSchema(schema interface{})
	SetResponseMIMEType(mimeType string)
	GetModelName() string
}

const (
	firstCatalogTemplate = `你是一位经验丰富的游戏客服对话分类专家。你的任务是将用户与客服的对话准确分类到预定义类别中。请仔细分析问题内容,考虑关键词和上下文,然后输出符合要求的JSON格式回复。

	<输入对话>
	%s
	</输入对话>
	
	<候选类别>
	充值问题: 包括各种支付平台未到账、充值失败、退款申请、恶意退款、商品信息获取问题、充值限制、第三方充值渠道相关问题等（通常表现为有“支付”或者“充值”这些关键词）。
	特殊问题: 包括游戏建议反馈、兑换码使用、预约奖励领取、道具兑换错误(如星轨专票/通票兑换错误)等非常规问题。
	法务相关问题: 涉及举报(如侵权、泄露、账号交易等)、数据披露等法律相关问题。
	HoYoLAB问题: 涉及HoYoLAB社区和工具(如崩坏：星穹铁道地图)使用的相关问题。
	账号问题: 包括账号安全(如被盗)、登录异常、信息修改(如改生日)、账号绑定解绑、忘记密码、验证码问题、账号注销、PSN相关问题等。
	活动问题: 主要涉及游戏内版本活动的相关问题,如活动奖励、活动任务、活动规则等（通常表现为有“活动/事件/event/quest"这些关键词）。
	游戏问题: 包括游戏系统(如跃迁、成就、模拟宇宙、背包、道具、商店、合成、养成功能等)、战斗及角色/怪物表现、任务流程(开拓、冒险、日常、活动任务等)、迷宫关卡(解密玩法、地图场景、物件表现等)、客户端问题(闪退、崩溃、显示异常、性能/卡顿、声音异常、手柄问题等)、游戏内本地化相关(文本翻译、配音等)、启动器问题等游戏内容和技术相关问题。
	</候选类别>
	
	<注意事项>
	- 活动问题优先级高于游戏问题，如果同时符合两者，则为活动问题。
	- 理由应简洁明了,引用原文关键词或短语
	- 如果无法确定分类,选择最相关的类别并详细说明原因
	- 结合整个问题的上下文来判断类别,不要仅依赖单个词语
	- 如果问题涉及多个方面,优先选择最主要或最核心的问题类别
	- 注意识别用户的具体诉求,而不仅仅是问题中提到的表面内容
	</注意事项>

	<输出格式>
	{"label": "游戏问题","reason": "用户反馈了一个关于游戏内遗物装备属性显示与实际效果不一致的问题。具体表现为'画像の遺物は速度差が\"2\"なのですが、装備すると速度が\"3\"変化します'，这属于游戏系统中的道具和角色属性表现问题，因此归类为游戏问题。"}
	</输出格式>
	`
	secondCatalogTemplate = `你是一位精通游戏行业的AI助手,专门负责分析玩家和客服对话内容，然后进行分类。你的任务是:
	1. 仔细阅读整个对话，理解玩家的问题和客服的回应。
	2. 识别对话中的关键词、短语和主题。
	3. 考虑游戏相关的上下文和常见问题类型。
	4. 从候选类别中选择最匹配的一级和二级子类别。
	5. 提供简洁但有力的理由，引用对话中的具体内容。

	<输入对话>
	%s
	</输入对话>

	<候选类别格式说明>
	候选类别是一个嵌套的JSON对象，结构如下：
	{
		"一级分类1": {
		"二级分类1": "二级分类1的描述",
		"二级分类2": "二级分类2的描述",
		...
		},
		"一级分类2": {
		"二级分类1": "二级分类1的描述",
		"二级分类2": "二级分类2的描述",
		...
		},
		...
	}
	每个一级子分类包含多个二级子分类，每个二级类别下都有其描述。
	在分类时，你只需要选择一级子分类和二级子分类。		
	</候选类别格式说明>

	<候选类别>
	%s
	</候选类别>

	<输出格式>
	{"primary_classification": "系统", "reason": "玩家反馈遗器装备后速度数值与描述不符，属于游戏内养成系统（遗器）的问题。", "secondary_classification": "养成功能"}
	</输出格式>


	注意事项：
	1. 确保选择的类别与对话内容高度相关。
	2. 理由应当简洁明了，直接引用对话中的关键内容。
	3. 如果对话涉及多个主题，请选择最主要或最重要的一个进行分类。
	4. 严格遵守JSON格式，正确处理特殊字符。
	5. 只输出JSON内容，不要添加任何其他解释或文字。
	6. 你需要使用转义字符，保证json的正确解析。	

	现在，请基于以上指导，直接输出符合要求的JSON格式回复。
`
)

func main() {
	file := flag.String("csv", "", "specify test csv file name")
	location := flag.String("location", "", "specify gcp location")
	modelName := flag.String("model", "", "specify gemini model name")
	projectID := flag.String("project", "", "specify gcp project id")
	maxConcurrent := flag.Int("max-concurrent", 5, "specify the max concurrent")
	flag.Parse()
	if *file == "" || *location == "" || *modelName == "" || *projectID == "" {
		log.Println("parameters can not be empty!")
		os.Exit(1)
	}
	log.Printf("location is %s\n", *location)
	log.Printf("modelName is %s\n", *modelName)
	log.Printf("projectID is %s\n", *projectID)

	var llm_api LLMAPI
	if strings.Contains(*modelName, "gemini") {
		llm_api = NewGeminiAPI(*location, *projectID, *modelName)
	} else {
		llm_api = NewClaudeAPI(*location, *projectID, *modelName)
	}

	records, err := readCSVToStruct(*file)
	if err != nil {
		log.Println("读取CSV文件失败:", err)
	}

	// 创建或打开输出文件
	outputFile, err := os.OpenFile("output.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	// 创建一个互斥锁，用于保护文件写入
	var fileMutex sync.Mutex

	var wg sync.WaitGroup
	// 创建带缓冲的 channel，容量为 maxConcurrent
	semaphore := make(chan struct{}, *maxConcurrent)

	for _, record := range records {
		wg.Add(1) // 每启动一个 goroutine，WaitGroup计数器加1
		go func(record Record) {
			defer wg.Done()         // goroutine 执行结束，WaitGroup计数器减1
			semaphore <- struct{}{} // 发送一个值到 channel，如果 channel 已满，则会阻塞
			defer func() { <-semaphore }()
			resp, err := getFirstLabel(record, llm_api)
			if err != nil {
				log.Printf("%s: %v", record.ID, err)
				return
			}
			// 使用互斥锁保护文件写入
			fileMutex.Lock()
			_, err = fmt.Fprintf(outputFile, "%s, %s-%s-%s, %s\n", resp.ID, resp.First_label, resp.Second_label, resp.Thrid_label, record.Label)
			if err != nil {
				log.Printf("%s写入文件失败: %v", record.ID, err)
			}
			fileMutex.Unlock()

		}(record) // 将 record 作为参数传递给 goroutine，防止闭包引用问题
	}
	wg.Wait() // 等待所有 goroutine 执行结束

}

func getFirstLabel(record Record, llm_api LLMAPI) (Result, error) {
	var first_invoked_response string

	if strings.Contains(llm_api.GetModelName(), "gemini") {
		llm_api.SetResponseMIMEType("application/json")
		llm_api.SetResponseSchema(&genai.Schema{Type: genai.TypeObject, Properties: map[string]*genai.Schema{
			"label":  {Type: genai.TypeString},
			"reason": {Type: genai.TypeString},
		}})
	}
	prompt := fmt.Sprintf(firstCatalogTemplate, record.Content)
	first_invoked_response, err := llm_api.InvokeText(prompt)
	if err != nil {
		return Result{}, err
	}
	// fmt.Println(first_invoked_response)

	var firstcatalog FirstCatalog
	json.Unmarshal([]byte(first_invoked_response), &firstcatalog)
	re := regexp.MustCompile(`[\n\\"]+`)
	firstcatalog.Label = re.ReplaceAllString(firstcatalog.Label, "")

	secondcatalog, err := getSecondLabel(firstcatalog.Label, llm_api, record.Content)
	if err != nil {
		return Result{}, err
	}
	var result Result
	result.ID = record.ID
	result.First_label = firstcatalog.Label
	secondcatalog.PrimaryClassification = re.ReplaceAllString(secondcatalog.PrimaryClassification, "")
	secondcatalog.SecondaryClassification = re.ReplaceAllString(secondcatalog.SecondaryClassification, "")
	result.Second_label = secondcatalog.PrimaryClassification
	result.Thrid_label = secondcatalog.SecondaryClassification
	result.Reason = secondcatalog.Reason
	return result, nil

}

func getSecondLabel(firstcatalog string, llm_api LLMAPI, content string) (SecondCatalog, error) {
	if strings.Contains(llm_api.GetModelName(), "gemini") {
		llm_api.SetResponseMIMEType("application/json")
		llm_api.SetResponseSchema(&genai.Schema{Type: genai.TypeObject, Properties: map[string]*genai.Schema{
			"primary_classification":   {Type: genai.TypeString},
			"secondary_classification": {Type: genai.TypeString},
			"reason":                   {Type: genai.TypeString},
		}})
	}

	second_label_rule := gjson.Get(string(string(labelRules)), firstcatalog)

	prompt := fmt.Sprintf(secondCatalogTemplate, content, second_label_rule)
	second_invoked_response, err := llm_api.InvokeText(prompt)
	// fmt.Println(second_invoked_response)
	if err != nil {
		return SecondCatalog{}, err
	}
	var secondcatalog SecondCatalog
	json.Unmarshal([]byte(second_invoked_response), &secondcatalog)
	return secondcatalog, nil
}

func readCSVToStruct(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("无法打开文件:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 跳过标题行
	_, err = reader.Read()
	if err != nil {
		log.Println("读取标题失败:", err)
		return nil, err
	}

	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("读取行错误:", err)
			return nil, err
		}
		record := Record{
			ID:      row[0],
			Result:  row[1],
			Comment: row[2],
			Label:   row[3],
			Content: row[4],
		}
		records = append(records, record)
	}
	return records, nil

}
