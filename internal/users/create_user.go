package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/GalahadKingsman/messenger_client/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) CreateUser(req models.CreateUserRequest) (*models.CreateUserResponse, error) {
	endpoint := fmt.Sprintf("%s/users/create", c.APIGatewayURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать тело запроса: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ: %w", err)
	}

	var gatewayResponse struct {
		Success string `json:"success"`
	}
	if err := json.Unmarshal(rawBody, &gatewayResponse); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	id, err := extractIDFromSuccessMessage(gatewayResponse.Success)
	if err != nil {
		return nil, fmt.Errorf("не удалось извлечь ID из ответа: %w", err)
	}

	return &models.CreateUserResponse{
		ID:      id,
		Success: gatewayResponse.Success,
	}, nil
}

// Вспомогательная функция для извлечения ID из строки
func extractIDFromSuccessMessage(msg string) (int64, error) {
	prefix := "Пользователь успешно создан с ID: "
	if !strings.HasPrefix(msg, prefix) {
		return 0, fmt.Errorf("неверный формат сообщения")
	}

	idStr := strings.TrimPrefix(msg, prefix)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать ID в число: %w", err)
	}

	return id, nil
}

func (u *UserCase) CreateUserCase() {
	var login, password, firstName, lastName, email, phone string

	_ = survey.AskOne(&survey.Input{Message: "Введите имя пользователя (login):"}, &login)
	_ = survey.AskOne(&survey.Input{Message: "Введите пароль:"}, &password)
	_ = survey.AskOne(&survey.Input{Message: "Введите имя:"}, &firstName)
	_ = survey.AskOne(&survey.Input{Message: "Введите фамилию:"}, &lastName)
	_ = survey.AskOne(&survey.Input{Message: "Введите email:"}, &email)
	_ = survey.AskOne(&survey.Input{Message: "Введите телефон:"}, &phone)

	req := models.CreateUserRequest{
		Login:     login,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	resp, err := u.client.CreateUser(req)
	if err != nil {
		fmt.Println("Ошибка при создании пользователя:", err)
		return
	}

	fmt.Println("Пользователь успешно создан, ID:", resp.ID)
}
