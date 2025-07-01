package dialog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io"
	"messenger_client/internal/models"
	"net/http"
)

func (c *Client) CreateDialog(ctx context.Context, req models.CreateDialogRequest, token string) (*models.CreateDialogResponse, error) {
	endpoint := fmt.Sprintf("%s/dialog/create", c.APIGatewayURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать запрос: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ: %w", err)
	}

	var result models.CreateDialogResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}

func (d *DialogCase) CreateDialogCase() {
	if d.token == "" {
		fmt.Println("Ошибка: необходимо сначала выполнить вход.")
		return
	}

	var userID, peerID int
	var dialogName string

	_ = survey.AskOne(&survey.Input{Message: "Введите свой UserID:"}, &userID)
	_ = survey.AskOne(&survey.Input{Message: "Введите PeerID собеседника:"}, &peerID)
	_ = survey.AskOne(&survey.Input{Message: "Введите название диалога:"}, &dialogName)

	req := models.CreateDialogRequest{
		UserID:     int32(userID),
		PeerID:     int32(peerID),
		DialogName: dialogName,
	}

	// вызываем обновлённый клиентский метод
	resp, err := d.client.CreateDialog(context.Background(), req, d.token)
	if err != nil {
		fmt.Println("Ошибка при создании диалога:", err)
		return
	}

	fmt.Println("Диалог успешно создан, ID:", resp.DialogID)
}
