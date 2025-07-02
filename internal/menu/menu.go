package menu

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"log"
	"messenger_client/internal/dialog"
	"messenger_client/internal/notifications"
	"messenger_client/internal/users"
	"os"
	"strings"
	"sync"
)

const apiURL = "http://127.0.0.1:8080"

func Run() {
	// Инициализация readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: ".messenger_history",
	})
	if err != nil {
		log.Fatalf("readline init error: %v", err)
	}
	defer rl.Close()

	// Инициализация клиентов и use-cases
	userClient := users.NewClient(apiURL)
	dlgClient := dialog.NewClient(apiURL)
	notifClient := notifications.NewClient(apiURL)

	userCase := users.NewUserCase(userClient)
	dlgCase := dialog.NewDialogCase(dlgClient)
	notifCase := notifications.NewNotCase(notifClient)

	for {
		fmt.Println("\nДоступные команды: Войти | Зарегистрироваться | Выход")
		line, err := rl.Readline()
		if err != nil {
			break
		}
		cmd := strings.TrimSpace(line)
		switch cmd {
		case "Войти":
			if userCase.LoginCase() {
				// Устанавливаем токен
				token := userCase.Token()
				dlgCase.SetToken(token)
				notifCase.SetToken(token)
				// Переходим в меню пользователя
				runUserMenu(rl, userCase, dlgCase, notifCase)
			}
		case "Зарегистрироваться":
			userCase.CreateUserCase()
		case "Выход":
			fmt.Println("Выход.")
			return
		default:
			fmt.Println("Неизвестная команда:", cmd)
		}
	}
}

func runUserMenu(
	rl *readline.Instance,
	userCase *users.UserCase,
	dlgCase *dialog.DialogCase,
	notifCase *notifications.NotificationCase,
) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Запускаем подписку на уведомления
	wg.Add(1)
	go notifCase.Listen(ctx, &wg)

	for {
		fmt.Println("\nДоступные команды: Получить пользователей | Создать диалог | Получить список диалогов | Отправить сообщение | Получить сообщения диалога | Выйти из аккаунта | Завершить приложение")
		line, err := rl.Readline()
		if err != nil {
			break
		}
		cmd := strings.TrimSpace(line)
		switch cmd {
		case "Получить пользователей":
			userCase.GetUsersCase()
		case "Создать диалог":
			dlgCase.CreateDialogCase()
		case "Получить список диалогов":
			dlgCase.GetUserDialogsCase()
		case "Отправить сообщение":
			dlgCase.SendMessageCase()
		case "Получить сообщения диалога":
			dlgCase.GetDialogMessagesCase()
		case "Выйти из аккаунта":
			cancel()
			wg.Wait()
			return
		case "Завершить приложение":
			cancel()
			wg.Wait()
			os.Exit(0)
		default:
			fmt.Println("Неизвестная команда:", cmd)
		}
	}
}
