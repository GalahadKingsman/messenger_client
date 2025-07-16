package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GalahadKingsman/messenger_client/internal/models"
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
			if resp.StatusCode == http.StatusNoContent {
				resp.Body.Close()
				time.Sleep(200 * time.Millisecond)
				continue
			}
			if resp.StatusCode == http.StatusOK {
				var notifs []models.Notification
				if err := json.NewDecoder(resp.Body).Decode(&notifs); err == nil {
					for _, n := range notifs {
						login := nc.getLogin(n.From)
						fmt.Printf("\n📩 Сообщение от %s, ID диалога – %d: %s\n",
							login,
							n.DialogID,
							n.Message,
						)
					}
				}
			}
			resp.Body.Close()
		}
	}
}

func (nc *NotificationCase) getLogin(userID string) string {
	url := fmt.Sprintf("%s/users/get?id=%s", nc.client.APIGatewayURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return userID
	}
	req.Header.Set("Authorization", "Bearer "+nc.token)

	resp, err := nc.client.HTTPClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return userID // fallback на ID
	}
	defer resp.Body.Close()

	var wrapper struct {
		Users []struct {
			Id    int64  `json:"id"`
			Login string `json:"login"`
		} `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return userID
	}
	if len(wrapper.Users) == 0 {
		return userID
	}
	return wrapper.Users[0].Login
}
