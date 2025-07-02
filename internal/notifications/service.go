package notifications

import (
	"net/http"
	"time"
)

type NotificationCase struct {
	client *Client
	token  string
}

type Client struct {
	APIGatewayURL string
	HTTPClient    *http.Client
}

func NewNotCase(client *Client) *NotificationCase {
	return &NotificationCase{client: client}
}

func NewClient(apiGatewayURL string) *Client {
	tr := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     60 * time.Second,
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	return &Client{
		APIGatewayURL: apiGatewayURL,
		HTTPClient:    httpClient,
	}
}

func (nc *NotificationCase) SetToken(token string) {
	nc.token = token
}
