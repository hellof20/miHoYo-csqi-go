# miHoYo-csqi-go

## Build
```
go build -o bin/csqi
```

## Run Gemini model
```
bin/csqi \
  -csv cs_qadata_fortest2.csv \
  -location us-central1 \
  -project speedy-victory-336109 \
  -model gemini-1.5-flash-002 \
  -max-concurrent 20
```
## Run Claude model
```
bin/csqi \
  -csv cs_qadata_fortest2.csv \
  -location us-east5 \
  -project speedy-victory-336109 \
  -model claude-3-5-sonnet@20240620 \
  -max-concurrent 10
```
