package users

import (
	"net/http"
	"time"
)

type UserCase struct {
	client *Client
	token  string
}

func NewUserCase(client *Client) *UserCase {
	return &UserCase{client: client}
}

type Client struct {
	APIGatewayURL string
	HTTPClient    *http.Client
}

func NewClient(apiGatewayURL string) *Client {
	tr := &http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}
	return &Client{
		APIGatewayURL: apiGatewayURL,
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		},
	}
}

func (s *UserCase) Token() string {
	return s.token
}
