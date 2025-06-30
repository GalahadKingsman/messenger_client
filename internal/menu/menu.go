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
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var tokenPath string

const apiURL = "http://localhost:8080"

func init() {
	home, _ := os.UserHomeDir()
	tokenPath = filepath.Join(home, ".messenger_token")
}

func Run() {
	apiURL := "http://localhost:8080"
	userClient := users.NewClient(apiURL)
	dialogClient := dialog.NewClient(apiURL)

	for {
		choice := ""
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: []string{"Войти", "Зарегистрироваться", "Выход"},
		}
		_ = survey.AskOne(prompt, &choice)

		switch choice {
		case "Войти":
			if users.LoginCase(userClient) {
				runUserMenu(userClient, dialogClient)
			}
		case "Зарегистрироваться":
			users.CreateUserCase(userClient)
		case "Выход":
			fmt.Println("Выход.")
			os.Exit(0)
		}
	}
}

func runUserMenu(userClient *users.Client, dialogClient *dialog.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenLoop(ctx, tokenPath)
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
			users.GetUsersCase(userClient)
		case "Создать диалог":
			dialog.CreateDialogCase(dialogClient)
		case "Получить список диалогов":
			dialog.GetUserDialogsCase(dialogClient)
		case "Отправить сообщение":
			dialog.SendMessageCase(dialogClient)
		case "Получить сообщения диалога":
			dialog.GetDialogMessagesCase(dialogClient)
		case "Выйти из аккаунта":
			os.Remove(tokenPath)
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

func listenLoop(ctx context.Context, tokenPath string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			token, err := os.ReadFile(tokenPath)
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
