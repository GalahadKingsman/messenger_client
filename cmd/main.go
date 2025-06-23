package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/dialog"
	"messenger_client/internal/models"
	"messenger_client/internal/users"
	"os"
	"strconv"
	"time"
)

const apiURL = "http://localhost:8080"

func main() {
	survey.WithStdio(os.Stdin, os.Stdout, os.Stderr)
	userClient := users.NewClient(apiURL)
	dialogClient := dialog.NewClient(apiURL)

	for {
		choice := ""
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: []string{
				"Создать пользователя",
				"Узнать свой ID",
				"Получить пользователей",
				"Создать диалог",
				"Получить список диалогов",
				"Отправить сообщение",
				"Получить сообщения диалога",
				"Выход",
			},
		}
		_ = survey.AskOne(prompt, &choice)

		switch choice {

		case "Создать пользователя":
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

			resp, err := userClient.CreateUser(req)
			if err != nil {
				fmt.Println("Ошибка при создании пользователя:", err)
				break
			}

			fmt.Println("Пользователь успешно создан, ID:", resp.Id)

		case "Узнать свой ID":
			var login, password string

			_ = survey.AskOne(&survey.Input{Message: "Введите имя пользователя (login):"}, &login)
			_ = survey.AskOne(&survey.Input{Message: "Введите пароль:"}, &password)

			req := models.LoginRequest{
				Login:    login,
				Password: password,
			}
			resp, err := userClient.Login(req)
			if err != nil {
				fmt.Println("Ошибка при входе:", err)
				break
			}
			fmt.Println("Ваш ID", resp.UserID)

		case "Получить пользователей":
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
				break
			}

			users, err := userClient.GetUsers(params)
			if err != nil {
				fmt.Println("Ошибка при получении пользователей:", err)
				break
			}

			fmt.Println("Пользователи:")
			for _, u := range users {
				fmt.Printf("Login: %s, Email: %s, Имя: %s %s, Телефон: %s\n",
					u.Login, u.Email, u.FirstName, u.LastName, u.Phone)
			}
		case "Создать диалог":
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
				break
			}

			fmt.Println("Диалог успешно создан, ID:", resp.DialogID)

		case "Получить список диалогов":
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
				break
			}

			fmt.Println("Диалоги пользователя:")
			for _, dlg := range resp.Dialogs {
				fmt.Printf("ID: %d, Последнее сообщение: %s\n", dlg.DialogID, dlg.LastMessage)
			}

		case "Отправить сообщение":
			var dialogID, userID int32
			var text string

			_ = survey.AskOne(&survey.Input{Message: "Введите ID диалога:"}, &dialogID)
			_ = survey.AskOne(&survey.Input{Message: "Введите свой UserID:"}, &userID)
			_ = survey.AskOne(&survey.Input{Message: "Введите текст сообщения:"}, &text)

			req := models.SendMessageRequest{
				DialogID: dialogID,
				UserID:   userID,
				Text:     text,
			}

			resp, err := dialogClient.SendMessage(req)
			if err != nil {
				fmt.Println("Ошибка при отправке сообщения:", err)
				break
			}

			fmt.Printf("Сообщение отправлено, ID: %d, Timestamp: %s\n", resp.MessageID, resp.Timestamp)

		case "Получить сообщения диалога":
			var dialogIDStr, limitStr, offsetStr string

			_ = survey.AskOne(&survey.Input{Message: "Введите ID диалога:"}, &dialogIDStr)
			_ = survey.AskOne(&survey.Input{Message: "Введите лимит сообщений (необязательно):"}, &limitStr)
			_ = survey.AskOne(&survey.Input{Message: "Введите, сколько сообщений пропустить (необязательно):"}, &offsetStr)

			dialogID64, err := strconv.ParseInt(dialogIDStr, 10, 32)
			if err != nil {
				fmt.Println("Некорректный ID диалога:", err)
				break
			}
			dialogID := int32(dialogID64)

			var limitPtr, offsetPtr *int32

			if limitStr != "" {
				limit64, err := strconv.ParseInt(limitStr, 10, 32)
				if err != nil {
					fmt.Println("Некорректный лимит:", err)
					break
				}
				limit := int32(limit64)
				limitPtr = &limit
			}

			if offsetStr != "" {
				offset64, err := strconv.ParseInt(offsetStr, 10, 32)
				if err != nil {
					fmt.Println("Некорректный offset:", err)
					break
				}
				offset := int32(offset64)
				offsetPtr = &offset
			}

			resp, err := dialogClient.GetDialogMessages(dialogID, limitPtr, offsetPtr)
			if err != nil {
				fmt.Println("Ошибка при получении сообщений диалога:", err)
				break
			}

			fmt.Println("Сообщения диалога:")
			for _, msg := range resp.Messages {
				fmt.Printf("ID: %d, UserID: %d, Text: %s, Timestamp: %s\n", msg.ID, msg.UserID, msg.Text, msg.Timestamp.Format(time.RFC822))
			}
		case "Выход":
			fmt.Println("Выход.")
			os.Exit(0)
		}
	}
}
