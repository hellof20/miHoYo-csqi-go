package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

type GeminiAPI struct {
	modelName        string
	projectID        string
	location         string
	responseSchema   *genai.Schema
	responseMIMEType string
}

func NewGeminiAPI(location, projectID, modelName string) *GeminiAPI {
	return &GeminiAPI{
		modelName: modelName,
		projectID: projectID,
		location:  location,
	}
}

func (a *GeminiAPI) GetModelName() string {
	return a.modelName
}

func (a *GeminiAPI) SetResponseSchema(schema interface{}) {
	a.responseSchema = schema.(*genai.Schema)
}

func (a *GeminiAPI) SetResponseMIMEType(mimeType string) {
	a.responseMIMEType = mimeType
}

func (a *GeminiAPI) InvokeText(prompt string) (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	client, err := genai.NewClient(ctx, a.projectID, a.location)
	if err != nil {
		return "", fmt.Errorf("new client with err: %w", err)
	}

	model := client.GenerativeModel(a.modelName)
	model.SetTemperature(0.1)
	if a.responseSchema != nil {
		model.GenerationConfig.ResponseSchema = a.responseSchema
		model.GenerationConfig.ResponseMIMEType = a.responseMIMEType
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", fmt.Errorf("generate content with err: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content")
	}

	result := resp.Candidates[0].Content.Parts[0]
	resultStr := fmt.Sprint(result)
	return resultStr, nil
}

func (a *GeminiAPI) InvokeImg(prompt string, img_path string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, a.projectID, a.location)
	if err != nil {
		fmt.Println(err)
	}

	model := client.GenerativeModel(a.modelName)
	model.SetTemperature(1)
	model.GenerationConfig.ResponseMIMEType = a.responseMIMEType
	model.GenerationConfig.ResponseSchema = a.responseSchema

	bytes, err := os.ReadFile(img_path)
	if err != nil {
		fmt.Println(err)
	}
	img_data := genai.ImageData("image/jpeg", bytes)
	prompt_data := genai.Text(prompt)

	resp, err := model.GenerateContent(ctx, img_data, prompt_data)
	if err != nil {
		fmt.Println(err)
	}

	result := resp.Candidates[0].Content.Parts[0]
	resultStr := fmt.Sprint(result)
	return resultStr, nil
}
