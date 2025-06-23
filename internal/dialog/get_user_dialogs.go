package dialog

import (
	"encoding/json"
	"fmt"
	"messenger_client/internal/models"
	"net/http"
)

func (c *Client) GetUserDialogs(userID int32, limit, offset *int32) (*models.GetUserDialogsResponse, error) {
	endpoint := fmt.Sprintf("%s/dialog/user?user_id=%d", c.APIGatewayURL, userID)

	if limit != nil {
		endpoint += fmt.Sprintf("&limit=%d", *limit)
	}
	if offset != nil {
		endpoint += fmt.Sprintf("&offset=%d", *offset)
	}

	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	var result models.GetUserDialogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}
