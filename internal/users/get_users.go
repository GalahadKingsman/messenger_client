package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"messenger_client/internal/models"
	"net/http"
	"net/url"
)

func (c *Client) GetUsers(params map[string]string) ([]models.CreateUserRequest, error) {
	endpoint := fmt.Sprintf("%s/users/get", c.APIGatewayURL)

	query := url.Values{}
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	if len(query) == 0 {
		return nil, errors.New("нужно указать хотя бы один параметр")
	}

	reqURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к API Gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	var result struct {
		Users []models.CreateUserRequest `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return result.Users, nil
}
