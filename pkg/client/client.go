package client

import (
	"fmt"
	"net/http"
)

type Client struct {
	AccessKey    string
	SecretKey    string
	SessionToken string

	http *http.Client
}

func NewClient(accessKey string, secretKey string, sessionToken string) *Client {
	return &Client{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
		http:         &http.Client{},
	}
}

func (c *Client) Do(req *http.Request, region string, service string, data []byte) (*http.Response, error) {
	req.Header.Set("Host", req.URL.Host)

	if len(data) > 0 {
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	}

	signRequest(req, c.AccessKey, c.SecretKey, c.SessionToken, region, service, data)

	return c.http.Do(req)
}
