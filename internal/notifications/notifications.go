package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Notification struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func ListenNotifications(ctx context.Context, wg *sync.WaitGroup, tokenPath string) {
	defer wg.Done()

	token, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Println("❌ Не удалось прочитать токен:", err)
		return
	}

	client := &http.Client{}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			req, _ := http.NewRequest("GET", "http://localhost:8080/notifications/longpoll", nil)
			req.Header.Set("Authorization", "Bearer "+string(token))

			resp, err := client.Do(req)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			if resp.StatusCode == http.StatusOK {
				var notifs []Notification
				if err := json.NewDecoder(resp.Body).Decode(&notifs); err == nil {
					for _, n := range notifs {
						fmt.Printf("\n🔔 %s: %s\n", n.From, n.Message)
					}
				}
			}
			resp.Body.Close()
		}
	}
}
