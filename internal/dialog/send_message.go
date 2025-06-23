package dialog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/models"
	"net/http"
)

func (c *Client) SendMessage(req models.SendMessageRequest) (*models.SendMessageResponse, error) {
	endpoint := fmt.Sprintf("%s/dialog/send", c.APIGatewayURL)

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

	var result models.SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}

func SendMessageCase(dialogClient *Client) {
	var dialogID, userID int32
	var text string

	_ = survey.AskOne(&survey.Input{Message: "Введите ID диалога:"}, &dialogID)
	_ = survey.AskOne(&survey.Input{Message: "Введите свой UserID:"}, &userID)
	_ = survey.AskOne(&survey.Input{Message: "Введите текст сообщения:"}, &text)

	req := models.SendMessageRequest{
		DialogID: dialogID,
		UserID:   userID,
		Text:     text,
	}

	resp, err := dialogClient.SendMessage(req)
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	fmt.Printf("Сообщение отправлено, ID: %d, Timestamp: %s\n", resp.MessageID, resp.Timestamp)
}
