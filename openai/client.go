package openai

import (
	"net/http"
)

type Client struct {
	requestConstructor *requestConstructor
	executor           *executor
	responseParser     *responseParser
}

type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		requestConstructor: &requestConstructor{APIKey: apiKey},
		executor:           &executor{client: &http.Client{}},
		responseParser:     &responseParser{},
	}
}

func (c *Client) ChatCompletion(prompt string) (*CompletionResponse, error) {
	req, err := c.requestConstructor.construct(prompt)
	if err != nil {
		return nil, err
	}

	respBody, err := c.executor.execute(req)
	if err != nil {
		return nil, err
	}

	return c.responseParser.parse(respBody)
}
