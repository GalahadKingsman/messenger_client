package menu

import (
	"context"
	"fmt"
	"github.com/GalahadKingsman/messenger_client/internal/dialog"
	"github.com/GalahadKingsman/messenger_client/internal/notifications"
	"github.com/GalahadKingsman/messenger_client/internal/users"
	"github.com/chzyer/readline"
	"log"
	"os"
	"strings"
	"sync"
)

const apiURL = "http://127.0.0.1:8080"

var mainCompleter = readline.NewPrefixCompleter(
	readline.PcItem("Войти"),
	readline.PcItem("Зарегистрироваться"),
	readline.PcItem("Выход"),
)

var userCompleter = readline.NewPrefixCompleter(
	readline.PcItem("Получить пользователей"),
	readline.PcItem("Создать диалог"),
	readline.PcItem("Получить список диалогов"),
	readline.PcItem("Отправить сообщение"),
	readline.PcItem("Получить сообщения диалога"),
	readline.PcItem("Выйти из аккаунта"),
	readline.PcItem("Завершить приложение"),
)

func Run() {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryFile:  ".messenger_history",
		AutoComplete: mainCompleter,
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

				token := userCase.Token()
				dlgCase.SetToken(token)
				notifCase.SetToken(token)
				rl.Config.AutoComplete = userCompleter

				runUserMenu(rl, userCase, dlgCase, notifCase)
				rl.Config.AutoComplete = mainCompleter
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
