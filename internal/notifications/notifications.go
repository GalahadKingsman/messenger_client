package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"messenger_client/internal/models"
	"net/http"
	"sync"
	"time"
)

func (nc *NotificationCase) Listen(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			req, _ := http.NewRequestWithContext(ctx,
				"GET",
				nc.client.APIGatewayURL+"/notifications/longpoll",
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+nc.token)

			resp, err := nc.client.HTTPClient.Do(req)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			if resp.StatusCode == http.StatusGatewayTimeout {
				resp.Body.Close()
				continue
			}
			if resp.StatusCode == http.StatusOK {
				var notifs []models.Notification
				if err := json.NewDecoder(resp.Body).Decode(&notifs); err == nil {
					for _, n := range notifs {
						fmt.Printf("\n📩 Сообщение от %s: %s\n", n.From, n.Message)
					}
				}
			}
			resp.Body.Close()
		}
	}
}
