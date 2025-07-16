package dialog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/GalahadKingsman/messenger_client/internal/models"
	"net/http"
)

func (c *Client) GetUserDialogs(
	ctx context.Context,
	userID int32,
	limit, offset *int32,
	token string,
) (*models.GetUserDialogsResponse, error) {
	endpoint := fmt.Sprintf("%s/dialog/user?user_id=%d", c.APIGatewayURL, userID)
	if limit != nil {
		endpoint += fmt.Sprintf("&limit=%d", *limit)
	}
	if offset != nil {
		endpoint += fmt.Sprintf("&offset=%d", *offset)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
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

func (d *DialogCase) GetUserDialogsCase() {
	if d.token == "" {
		fmt.Println("Ошибка: необходимо сначала выполнить вход.")
		return
	}

	var userID int32
	var limit, offset int32
	var applyLimit, applyOffset bool

	_ = survey.AskOne(&survey.Input{Message: "Введите ID пользователя:"}, &userID)
	_ = survey.AskOne(&survey.Confirm{
		Message: "Хотите указать лимит?",
		Default: false,
	}, &applyLimit)
	if applyLimit {
		_ = survey.AskOne(&survey.Input{Message: "Введите лимит:"}, &limit)
	}

	_ = survey.AskOne(&survey.Confirm{
		Message: "Хотите указать, сколько диалогов пропустить?",
		Default: false,
	}, &applyOffset)
	if applyOffset {
		_ = survey.AskOne(&survey.Input{Message: "Введите количество для смещения:"}, &offset)
	}

	var limitPtr, offsetPtr *int32
	if applyLimit {
		limitPtr = &limit
	}
	if applyOffset {
		offsetPtr = &offset
	}

	resp, err := d.client.GetUserDialogs(
		context.Background(),
		userID,
		limitPtr,
		offsetPtr,
		d.token,
	)
	if err != nil {
		fmt.Println("Ошибка при получении списка диалогов:", err)
		return
	}

	fmt.Println("Диалоги пользователя:")
	for _, dlg := range resp.Dialogs {
		fmt.Printf("ID: %d, Последнее сообщение: %s\n", dlg.DialogID, dlg.LastMessage)
	}
}
