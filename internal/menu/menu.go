package menu

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/dialog"
	"messenger_client/internal/users"
	"os"
)

const apiURL = "http://localhost:8080"

func Run() {
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
			users.CreateUserCase(userClient)

		case "Узнать свой ID":
			users.LoginCase(userClient)

		case "Получить пользователей":
			users.GetUsersCase(userClient)

		case "Создать диалог":
			dialog.CreateDialogCase(dialogClient)

		case "Получить список диалогов":
			dialog.GetUserDialogsCase(dialogClient)

		case "Отправить сообщение":
			dialog.SendMessageCase(dialogClient)

		case "Получить сообщения диалога":
			dialog.GetDialogMessagesCase(dialogClient)
		case "Выход":
			fmt.Println("Выход.")
			os.Exit(0)
		}
	}
}
