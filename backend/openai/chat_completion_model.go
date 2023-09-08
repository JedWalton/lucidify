package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tiktoken-go/tokenizer"
)

const baseURL = "https://api.openai.com/v1"

type Client struct {
	APIKey             string
	session            ChatSession
	requestConstructor requestConstructorInterface
	executor           executorInterface
	responseParser     responseParserInterface
}

type requestConstructorInterface interface {
	construct(messages []Message) (*http.Request, error)
}

type executorInterface interface {
	execute(req *http.Request) ([]byte, error)
}

type responseParserInterface interface {
	parse(body []byte) (*CompletionResponse, error)
}

type chatCompletionPayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type requestConstructor struct {
	APIKey string
}

type ChatSession struct {
	Messages []Message
}

type executor struct {
	client *http.Client
}

type responseParser struct{}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:             apiKey,
		requestConstructor: &requestConstructor{APIKey: apiKey},
		executor:           &executor{client: &http.Client{}},
		responseParser:     &responseParser{},
	}
}

// Send a message and maintain the context
func (c *Client) SendMessage(userInput string, systemInput string) (*CompletionResponse, error) {
	c.session.AddMessage("user", userInput)
	c.session.AddMessage("system", systemInput)

	req, err := c.requestConstructor.construct(c.session.Messages)
	if err != nil {
		return nil, err
	}

	respBody, err := c.executor.execute(req)
	if err != nil {
		return nil, err
	}

	response, err := c.responseParser.parse(respBody)
	if err != nil {
		return nil, err
	}

	// Add the assistant's response to the session
	c.session.AddMessage("assistant", response.Choices[0].Message.Content)

	return response, nil
}

func (s *ChatSession) AddMessage(role string, content string) {
	s.Messages = append(s.Messages, Message{Role: role, Content: content})
}

func EstimateTokenCount(text string) int {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		panic("Failed to get tokenizer: " + err.Error())
	}

	// Encode the input string to get a list of token ids
	ids, _, err := enc.Encode(text)
	if err != nil {
		panic("Failed to encode string: " + err.Error())
	}
	// Count and print the number of tokens
	tokenCount := len(ids)

	return tokenCount
}

func EstimateTokenCountOfCurrentChatSession(chatSession ChatSession) int {
	var tokenCount int
	for _, message := range chatSession.Messages {
		tokenCount += EstimateTokenCount(message.Content)
	}
	return tokenCount
}

func (rc *requestConstructor) construct(messages []Message) (*http.Request, error) {
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	payload := chatCompletionPayload{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}
	const maxPayloadSize = 2e6 // For instance, 2MB

	if len(data) > maxPayloadSize {
		return nil, fmt.Errorf("payload too large")
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
