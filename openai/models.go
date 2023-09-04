package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.openai.com/v1"

type chatCompletionPayload struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestConstructor struct {
	APIKey string
}

type responseParser struct{}

type executor struct {
	client *http.Client
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

func (rc *requestConstructor) construct(prompt string) (*http.Request, error) {
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	payload := chatCompletionPayload{
		Model: "gpt-3.5-turbo",
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rc.APIKey))

	return req, nil
}

func (e *executor) execute(req *http.Request) ([]byte, error) {
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func (rp *responseParser) parse(body []byte) (*CompletionResponse, error) {
	var completionResp CompletionResponse
	err := json.Unmarshal(body, &completionResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &completionResp, nil
}
