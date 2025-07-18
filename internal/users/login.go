package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/GalahadKingsman/messenger_client/internal/models"
	"net/http"
)

func (c *Client) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	endpoint := fmt.Sprintf("%s/users/login", c.APIGatewayURL)

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

	var result models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return &result, nil
}

func (u *UserCase) LoginCase() bool {
	var login, password string

	_ = survey.AskOne(&survey.Input{Message: "Введите имя пользователя:"}, &login)
	_ = survey.AskOne(&survey.Password{Message: "Введите пароль:"}, &password)

	req := models.LoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := u.client.Login(req)
	if err != nil {
		fmt.Println("Ошибка при входе:", err)
		return false
	}

	if resp.Token == "" {
		fmt.Println("Неверные данные")
		return false
	}

	u.token = resp.Token

	fmt.Println("Успешный вход. Ваш ID:", resp.UserID)
	return true
}
