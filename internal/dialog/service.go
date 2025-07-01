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
	return &Client{
		APIGatewayURL: apiGatewayURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (d *DialogCase) SetToken(tok string) {
	d.token = tok
}
