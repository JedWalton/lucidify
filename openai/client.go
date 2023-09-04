package openai

import (
	"net/http"
)

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
