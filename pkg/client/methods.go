package client

import (
	"net/http"
)

func (c *Client) Send(method string, url string, body []uint8) (*ClientResponse, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &ClientResponse{ResponseCode: resp.StatusCode}, nil
}
