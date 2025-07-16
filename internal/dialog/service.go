package dialog

import (
	"net/http"
	"time"
)

type DialogCase struct {
	client *Client
	token  string
}

func NewDialogCase(client *Client) *DialogCase {
	return &DialogCase{client: client}
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

func (d *DialogCase) SetToken(tok string) {
	d.token = tok
}
