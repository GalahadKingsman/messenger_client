package dialog

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
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

func GetUserDialogsCase(dialogClient *Client) {
	var userID int32
	var limit, offset int32
	var applyLimit, applyOffset bool

	_ = survey.AskOne(&survey.Input{Message: "Введите ID пользователя:"}, &userID)
	_ = survey.AskOne(&survey.Confirm{Message: "Хотите указать лимит?", Default: false}, &applyLimit)
	if applyLimit {
		_ = survey.AskOne(&survey.Input{Message: "Введите лимит:"}, &limit)
	}

	_ = survey.AskOne(&survey.Confirm{Message: "Хотите указать, сколько диалогов пропустить?", Default: false}, &applyOffset)
	if applyOffset {
		_ = survey.AskOne(&survey.Input{Message: "Введите, сколько диалогов пропустить:"}, &offset)
	}

	var limitPtr, offsetPtr *int32
	if applyLimit {
		limitPtr = &limit
	}
	if applyOffset {
		offsetPtr = &offset
	}

	resp, err := dialogClient.GetUserDialogs(userID, limitPtr, offsetPtr)
	if err != nil {
		fmt.Println("Ошибка при получении списка диалогов:", err)
		return
	}

	fmt.Println("Диалоги пользователя:")
	for _, dlg := range resp.Dialogs {
		fmt.Printf("ID: %d, Последнее сообщение: %s\n", dlg.DialogID, dlg.LastMessage)
	}
}
