package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/models"
	"net/http"
	"net/url"
)

func (c *Client) GetUsers(params map[string]string) ([]models.CreateUserRequest, error) {
	endpoint := fmt.Sprintf("%s/users/get", c.APIGatewayURL)

	query := url.Values{}
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	if len(query) == 0 {
		return nil, errors.New("нужно указать хотя бы один параметр")
	}

	reqURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к API Gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	var result struct {
		Users []models.CreateUserRequest `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return result.Users, nil
}

func GetUsersCase(userClient *Client) {
	params := make(map[string]string)

	var login, firstName, lastName, email, phone string

	_ = survey.AskOne(&survey.Input{Message: "Фильтр по логину (оставьте пустым, если не нужен):"}, &login)
	_ = survey.AskOne(&survey.Input{Message: "Фильтр по имени (оставьте пустым, если не нужен):"}, &firstName)
	_ = survey.AskOne(&survey.Input{Message: "Фильтр по фамилии (оставьте пустым, если не нужен):"}, &lastName)
	_ = survey.AskOne(&survey.Input{Message: "Фильтр по email (оставьте пустым, если не нужен):"}, &email)
	_ = survey.AskOne(&survey.Input{Message: "Фильтр по телефону (оставьте пустым, если не нужен):"}, &phone)

	if login != "" {
		params["login"] = login
	}
	if firstName != "" {
		params["first_name"] = firstName
	}
	if lastName != "" {
		params["last_name"] = lastName
	}
	if email != "" {
		params["email"] = email
	}
	if phone != "" {
		params["phone"] = phone
	}

	if len(params) == 0 {
		fmt.Println("Нужно указать хотя бы один параметр фильтра.")
		return
	}

	users, err := userClient.GetUsers(params)
	if err != nil {
		fmt.Println("Ошибка при получении пользователей:", err)
		return
	}

	fmt.Println("Пользователи:")
	for _, u := range users {
		fmt.Printf("Login: %s, Email: %s, Имя: %s %s, Телефон: %s\n",
			u.Login, u.Email, u.FirstName, u.LastName, u.Phone)
	}
}
