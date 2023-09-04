package openai

import (
	"net/http"
)

type requestConstructor struct {
	APIKey string
}

type responseParser struct{}

type executor struct {
	client *http.Client
}

type Client struct {
	requestConstructor *requestConstructor
	executor           *executor
	responseParser     *responseParser
}

func NewClient(apiKey string) *Client {
	return &Client{
		requestConstructor: &requestConstructor{APIKey: apiKey},
		executor:           &executor{client: &http.Client{}},
		responseParser:     &responseParser{},
	}
}
