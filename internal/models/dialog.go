package models

import "time"

type CreateDialogRequest struct {
	UserID     int32  `json:"user_id"`
	PeerID     int32  `json:"peer_id"`
	DialogName string `json:"dialog_name"`
}

type CreateDialogResponse struct {
	DialogID   int32  `json:"dialog_id"`
	DialogName string `json:"dialog_name"`
	Success    bool   `json:"success"`
}

type Message struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type GetDialogMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type SendMessageRequest struct {
	DialogID int32  `json:"dialog_id"`
	UserID   int32  `json:"user_id"`
	Text     string `json:"text"`
}

type SendMessageResponse struct {
	MessageID int64  `json:"message_id"`
	Timestamp string `json:"timestamp"`
}

type Dialog struct {
	DialogID    int32  `json:"dialog_id"`
	PeerID      int32  `json:"peer_id"`
	PeerLogin   string `json:"peer_login"`
	LastMessage string `json:"last_message"`
}

type GetUserDialogsResponse struct {
	Dialogs []Dialog `json:"dialogs"`
}
