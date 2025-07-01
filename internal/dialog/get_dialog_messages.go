package dialog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (c *Client) GetDialogMessages(
	ctx context.Context,
	dialogID int32,
	limit, offset *int32,
	token string,
) (*models.GetDialogMessagesResponse, error) {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	// добавляем токен в заголовок
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
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

func (d *DialogCase) GetDialogMessagesCase() {
	if d.token == "" {
		fmt.Println("Ошибка: необходимо сначала выполнить вход.")
		return
	}

	var dialogIDStr, limitStr, offsetStr string

	_ = survey.AskOne(&survey.Input{Message: "Введите ID диалога:"}, &dialogIDStr)
	_ = survey.AskOne(&survey.Input{Message: "Введите лимит сообщений (необязательно):"}, &limitStr)
	_ = survey.AskOne(&survey.Input{Message: "Введите, сколько сообщений пропустить (необязательно):"}, &offsetStr)

	// парсим dialogID
	dialogID64, err := strconv.ParseInt(dialogIDStr, 10, 32)
	if err != nil {
		fmt.Println("Некорректный ID диалога:", err)
		return
	}
	dialogID := int32(dialogID64)

	// парсим limit и offset
	var limitPtr, offsetPtr *int32
	if limitStr != "" {
		v, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			fmt.Println("Некорректный лимит:", err)
			return
		}
		tmp := int32(v)
		limitPtr = &tmp
	}
	if offsetStr != "" {
		v, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			fmt.Println("Некорректный offset:", err)
			return
		}
		tmp := int32(v)
		offsetPtr = &tmp
	}

	// вызываем обновлённый клиентский метод
	resp, err := d.client.GetDialogMessages(
		context.Background(),
		dialogID,
		limitPtr,
		offsetPtr,
		d.token,
	)
	if err != nil {
		fmt.Println("Ошибка при получении сообщений диалога:", err)
		return
	}

	// выводим результат
	fmt.Println("Сообщения диалога:")
	for _, msg := range resp.Messages {
		// msg.Timestamp — time.Time?
		fmt.Printf(
			"ID: %d, UserID: %d, Text: %s, Timestamp: %s\n",
			msg.ID,
			msg.UserID,
			msg.Text,
			msg.Timestamp.Format(time.RFC822),
		)
	}
}
