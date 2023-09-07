package client

import (
	"log"
	"net/http"
)

func (c *Client) Get(url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("fail in send postback, error - %s", err)
		return
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("fail in send postback, error - %s", err)
		return
	}

	log.Printf("response status code - %d", resp.StatusCode)
}
