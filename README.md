# miHoYo-csqi-go

## Build
```
go build -o bin/csqi
```

## Run Gemini model
```
bin/csqi \
  -csv data/cs_qadata_fortest2.csv \
  -location us-central1 \
  -project speedy-victory-336109 \
  -model gemini-1.5-flash-002 \
  -max-concurrent 20
```
## Run Claude model
```
bin/csqi \
  -csv data/cs_qadata_fortest2.csv \
  -location us-east5 \
  -project speedy-victory-336109 \
  -model claude-3-5-sonnet@20240620 \
  -max-concurrent 10
```

## Result
程序运行后会生成output.csv，第一个字段为ID，第二个字段为大模型生成的客诉类别，第三个字段为人工的客诉类别

<img width="955" alt="image" src="https://github.com/user-attachments/assets/8ea67058-0c4c-4e4d-bf75-398f451bafd8" />
