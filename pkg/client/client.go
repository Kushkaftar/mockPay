package client

import "net/http"

type PostbackStatusCode struct {
	StatusCode int
}

type Client struct {
	HttpClient http.Client
}

func NewClient() *Client {

	return &Client{
		HttpClient: http.Client{},
	}
}
