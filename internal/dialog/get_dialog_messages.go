package dialog

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

func GetDialogMessagesCase(dialogClient *Client) {
	var dialogIDStr, limitStr, offsetStr string

	_ = survey.AskOne(&survey.Input{Message: "Введите ID диалога:"}, &dialogIDStr)
	_ = survey.AskOne(&survey.Input{Message: "Введите лимит сообщений (необязательно):"}, &limitStr)
	_ = survey.AskOne(&survey.Input{Message: "Введите, сколько сообщений пропустить (необязательно):"}, &offsetStr)

	dialogID64, err := strconv.ParseInt(dialogIDStr, 10, 32)
	if err != nil {
		fmt.Println("Некорректный ID диалога:", err)
		return
	}
	dialogID := int32(dialogID64)

	var limitPtr, offsetPtr *int32

	if limitStr != "" {
		limit64, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			fmt.Println("Некорректный лимит:", err)
			return
		}
		limit := int32(limit64)
		limitPtr = &limit
	}

	if offsetStr != "" {
		offset64, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			fmt.Println("Некорректный offset:", err)
			return
		}
		offset := int32(offset64)
		offsetPtr = &offset
	}

	resp, err := dialogClient.GetDialogMessages(dialogID, limitPtr, offsetPtr)
	if err != nil {
		fmt.Println("Ошибка при получении сообщений диалога:", err)
		return
	}

	fmt.Println("Сообщения диалога:")
	for _, msg := range resp.Messages {
		fmt.Printf("ID: %d, UserID: %d, Text: %s, Timestamp: %s\n", msg.ID, msg.UserID, msg.Text, msg.Timestamp.Format(time.RFC822))
	}
}
