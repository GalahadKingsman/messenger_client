package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"messenger_client/internal/dialog"
	"messenger_client/internal/notifications"
	"messenger_client/internal/users"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)


const apiURL = "http://localhost:8080"

func Run() {
	baseClientUser := users.NewClient(apiURL)
	baseClientDialog := dialog.NewClient(apiURL)
	userCase := users.NewUserCase(baseClientUser)
	dialogCase := dialog.NewDialogCase(baseClientDialog)

	for {
		choice := ""
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: []string{"Войти", "Зарегистрироваться", "Выход"},
		}
		_ = survey.AskOne(prompt, &choice)

		switch choice {
		case "Войти":
			if userCase.LoginCase() {
				dialogCase.SetToken(userCase.Token())
				runUserMenu(userCase, dialogCase)
			}
		case "Зарегистрироваться":
			userCase.CreateUserCase()
		case "Выход":
			fmt.Println("Выход.")
			os.Exit(0)
		}
	}
}

func runUserMenu(userCase *users.UserCase, dialogCase *dialog.DialogCase) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenLoop(ctx, userCase.Token())
	}()

	for {
		choice := ""
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: []string{
				"Получить пользователей",
				"Создать диалог",
				"Получить список диалогов",
				"Отправить сообщение",
				"Получить сообщения диалога",
				"Выйти из аккаунта",
				"Завершить приложение",
			},
		}
		_ = survey.AskOne(prompt, &choice)

		switch choice {
		case "Получить пользователей":
			userCase.GetUsersCase()
		case "Создать диалог":
			dialogCase.CreateDialogCase()
		case "Получить список диалогов":
			dialogCase.GetUserDialogsCase()
		case "Отправить сообщение":
			dialogCase.SendMessageCase()
		case "Получить сообщения диалога":
			dialogCase.GetDialogMessagesCase()
		case "Выйти из аккаунта":
			cancel()
			wg.Wait()
			return
		case "Завершить приложение":
			cancel()
			wg.Wait()
			os.Exit(0)
		}
	}
}

func listenLoop(ctx context.Context, token string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			token, err := os.ReadFile(token)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8082/notifications/longpoll", nil)
			req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			if resp.StatusCode == http.StatusGatewayTimeout {
				resp.Body.Close()
				continue // снова ждать
			}

			var notifs []notifications.Notification
			err = json.NewDecoder(resp.Body).Decode(&notifs)
			resp.Body.Close()
			if err == nil {
				for _, n := range notifs {
					fmt.Printf("\n📩 Сообщение от %s: %s\n", n.From, n.Message)
				}
			}
		}
	}
}
