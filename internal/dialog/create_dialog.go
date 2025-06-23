package dialog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"messenger_client/internal/models"
	"net/http"
	"time"
)

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

func (c *Client) CreateDialog(req models.CreateDialogRequest) (*models.CreateDialogResponse, error) {
	endpoint := fmt.Sprintf("%s/dialog/create", c.APIGatewayURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать запрос: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	var result models.CreateDialogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}
