package dialog

import (
	"encoding/json"
	"fmt"
	"messenger_client/internal/models"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) GetDialogMessages(dialogID int32, limit *int32, offset *int32) (*models.GetDialogMessagesResponse, error) {
	baseURL := fmt.Sprintf("%s/dialog/messages", c.APIGatewayURL)

	params := url.Values{}
	params.Set("dialog_id", strconv.Itoa(int(dialogID)))
	if limit != nil {
		params.Set("limit", strconv.Itoa(int(*limit)))
	}
	if offset != nil {
		params.Set("offset", strconv.Itoa(int(*offset)))
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	var result models.GetDialogMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}
