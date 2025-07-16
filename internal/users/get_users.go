package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/GalahadKingsman/messenger_client/internal/models"
	"net/http"
)

func (c *Client) GetUsers(
	ctx context.Context,
	params map[string]string,
	token string,
) ([]models.UserResponse, error) {
	// 1) Собираем URL
	endpoint := fmt.Sprintf("%s/users/get", c.APIGatewayURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	// 2) Добавляем все фильтры в query
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// 3) Добавляем заголовок
	req.Header.Set("Authorization", "Bearer "+token)

	// 4) Выполняем запрос
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Gateway вернул статус %d", resp.StatusCode)
	}

	// 5) Десериализуем обёртку
	var wrapper struct {
		Users []models.UserResponse `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ: %w", err)
	}

	return wrapper.Users, nil
}

func (u *UserCase) GetUsersCase() {
	// 1) Проверка, что мы залогинены
	if u.token == "" {
		fmt.Println("Ошибка: необходимо сначала выполнить вход.")
		return
	}

	// 2) Собираем фильтры
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

	// 3) Делаем запрос к API, передавая токен
	ctx := context.Background()
	usersList, err := u.client.GetUsers(ctx, params, u.token)
	if err != nil {
		fmt.Println("Ошибка при получении пользователей:", err)
		return
	}

	// 4) Выводим результат
	fmt.Println("Пользователи:")
	for _, usr := range usersList {
		fmt.Printf("ID:%d Login:%s Имя:%s %s Email:%s Телефон:%s\n",
			usr.ID, usr.Login, usr.FirstName, usr.LastName, usr.Email, usr.Phone)
	}
}
