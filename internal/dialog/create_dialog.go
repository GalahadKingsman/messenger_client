package dialog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/models"
	"net/http"
	"time"
)

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

func CreateDialogCase(dialogClient *Client) {
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

	resp, err := dialogClient.CreateDialog(req)
	if err != nil {
		fmt.Println("Ошибка при создании диалога:", err)
		return
	}

	fmt.Println("Диалог успешно создан, ID:", resp.DialogID)
}
