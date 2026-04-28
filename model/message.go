package model

import (
	"encoding/json/jsontext"
	"time"

	"github.com/google/uuid"
)

// type ChatType string

// const (
// 	ChatTypeSingle ChatType = "single"
// 	ChatTypeGroup  ChatType = "group"
// )

// type MessageType string

// const (
// 	MessageTypeText  MessageType = "text"
// 	MessageTypeImage MessageType = "image"
// 	MessageTypeVideo MessageType = "video"
// 	MessageTypeFile  MessageType = "file"
// )

type JsonChatMessage struct {
	MsgID        uuid.UUID      `json:"msg_id"`
	ClientMsgID  uuid.UUID      `json:"client_msg_id"`
	SenderID     uuid.UUID      `json:"sender_id"`
	ReceiverID   uuid.UUID      `json:"receiver_id"`
	ChatType     ChatType       `json:"chat_type"`
	ServerTime   int64          `json:"server_time"`
	ReplyToMsgID *uuid.UUID     `json:"reply_to_msg_id"`
	MsgType      MessageType    `json:"msg_type"`
	Payload      jsontext.Value `json:"payload"`
	Ext          jsontext.Value `json:"ext"`
	CreatedAt    time.Time      `json:"created_at"`
}
