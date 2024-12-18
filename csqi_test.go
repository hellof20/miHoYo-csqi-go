package main

import (
	"testing"
)

func TestGemini(t *testing.T) {
	api := NewGeminiAPI("us-central1", "speedy-victory-336109", "gemini-1.5-flash-002")
	resp, err := api.InvokeText("hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp=%v", resp)
}

func TestClaude(t *testing.T) {
	api := NewClaudeAPI("us-east5", "speedy-victory-336109", "claude-3-5-sonnet@20240620")
	prompt := `
	2024-09-20 13:54:20 [玩家] ***:
    
    問題詳細: 画像の遺物は速度差が""2""なのですが、装備すると速度が""3""変化します。不具合でしょうか？
    画像アップロード（最大アップロード数は5枚、10MB以下）: 
        <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/******df81e043ff******_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/712fff53c7ac5ab6d4a4a14e5183b2ce_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/ca5306b45fddcd3edf9232a76fe7a679_******.png"" style=""max-height: 150px;"" ></img>
    動画アップロード （最大アップロード数は10本、各動画のサイズは200MB以下、.mp4、 .webm、 .ogg、.movのみサポート）: 
    Email: ***@***.***

2024-09-20 15:04:28 [客服] ***:
    親愛なる開拓者様
    この度は「崩壊：スターレイル」へのお問い合わせ、誠にありがとうございます。
    お問い合わせの件に関しまして、ご報告いただき誠にありがとうございます。
    ご申告事象の確認の結果、再現性を確認し、問題がございました際には修正を進めてまいりたく存じます。
    なお、仕様によるものの場合、現状のままとなりますことご了承いただきますようお願いいたします。
    また攻略情報に繋がることがございますため、仕様面についてはご質問にお答えいたしかねてしまうケースもございますこと、
    何卒ご了承いただけますと幸いです。
    この度は、ご連絡をいただきましたこと心よりお礼申し上げます。
    開拓者様がより楽しめるゲームとなるよう努めてまいる所存です。
    引き続き「崩壊：スターレイル」をご愛顧いただきますよう
    よろしくお願いいたします。
    Regards,
    カスタマーサポート担当
    ※本メールの内容の無断掲載、無断複製、転送は禁止させていただいております。ご注意ください。`
	resp, err := api.InvokeText(prompt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp=%v", resp)
}

func TestGeminiGetLabel(t *testing.T) {
	api := LLMAPI(NewGeminiAPI("us-central1", "speedy-victory-336109", "gemini-1.5-flash-002"))
	record := Record{
		ID:      "115884368906158080",
		Result:  "",
		Comment: "",
		Label:   "",
		Content: `2024-09-20 13:54:20 [玩家] ***:
    
    問題詳細: 画像の遺物は速度差が""2""なのですが、装備すると速度が""3""変化します。不具合でしょうか？
    画像アップロード（最大アップロード数は5枚、10MB以下）: 
        <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/******df81e043ff******_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/712fff53c7ac5ab6d4a4a14e5183b2ce_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/ca5306b45fddcd3edf9232a76fe7a679_******.png"" style=""max-height: 150px;"" ></img>
    動画アップロード （最大アップロード数は10本、各動画のサイズは200MB以下、.mp4、 .webm、 .ogg、.movのみサポート）: 
    Email: ***@***.***

2024-09-20 15:04:28 [客服] ***:
    親愛なる開拓者様
    この度は「崩壊：スターレイル」へのお問い合わせ、誠にありがとうございます。
    お問い合わせの件に関しまして、ご報告いただき誠にありがとうございます。
    ご申告事象の確認の結果、再現性を確認し、問題がございました際には修正を進めてまいりたく存じます。
    なお、仕様によるものの場合、現状のままとなりますことご了承いただきますようお願いいたします。
    また攻略情報に繋がることがございますため、仕様面についてはご質問にお答えいたしかねてしまうケースもございますこと、
    何卒ご了承いただけますと幸いです。
    この度は、ご連絡をいただきましたこと心よりお礼申し上げます。
    開拓者様がより楽しめるゲームとなるよう努めてまいる所存です。
    引き続き「崩壊：スターレイル」をご愛顧いただきますよう
    よろしくお願いいたします。
    Regards,
    カスタマーサポート担当
    ※本メールの内容の無断掲載、無断複製、転送は禁止させていただいております。ご注意ください。
`,
	}
	resp, err := getFirstLabel(record, api)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp=%v", resp)
}

func TestClaudeGetLabel(t *testing.T) {
	api := LLMAPI(NewClaudeAPI("us-east5", "speedy-victory-336109", "claude-3-5-sonnet@20240620"))
	record := Record{
		ID:      "115884368906158080",
		Result:  "",
		Comment: "",
		Label:   "",
		Content: `2024-09-20 13:54:20 [玩家] ***:
    
    問題詳細: 画像の遺物は速度差が""2""なのですが、装備すると速度が""3""変化します。不具合でしょうか？
    画像アップロード（最大アップロード数は5枚、10MB以下）: 
        <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/******df81e043ff******_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/712fff53c7ac5ab6d4a4a14e5183b2ce_******.png"" style=""max-height: 150px;"" ></img> <img src=""https://operation-private.inc-static.hoyoverse.com/csc-user/2024/09/20/0/ca5306b45fddcd3edf9232a76fe7a679_******.png"" style=""max-height: 150px;"" ></img>
    動画アップロード （最大アップロード数は10本、各動画のサイズは200MB以下、.mp4、 .webm、 .ogg、.movのみサポート）: 
    Email: ***@***.***

2024-09-20 15:04:28 [客服] ***:
    親愛なる開拓者様
    この度は「崩壊：スターレイル」へのお問い合わせ、誠にありがとうございます。
    お問い合わせの件に関しまして、ご報告いただき誠にありがとうございます。
    ご申告事象の確認の結果、再現性を確認し、問題がございました際には修正を進めてまいりたく存じます。
    なお、仕様によるものの場合、現状のままとなりますことご了承いただきますようお願いいたします。
    また攻略情報に繋がることがございますため、仕様面についてはご質問にお答えいたしかねてしまうケースもございますこと、
    何卒ご了承いただけますと幸いです。
    この度は、ご連絡をいただきましたこと心よりお礼申し上げます。
    開拓者様がより楽しめるゲームとなるよう努めてまいる所存です。
    引き続き「崩壊：スターレイル」をご愛顧いただきますよう
    よろしくお願いいたします。
    Regards,
    カスタマーサポート担当
    ※本メールの内容の無断掲載、無断複製、転送は禁止させていただいております。ご注意ください。
`,
	}
	resp, err := getFirstLabel(record, api)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp=%v", resp)
}
