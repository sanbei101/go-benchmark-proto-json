package model

import (
	"encoding/json/jsontext"

	"github.com/google/uuid"
)

type JsonChatMessage struct {
	MsgID        uuid.UUID      `json:"msg_id"`
	ClientMsgID  uuid.UUID      `json:"client_msg_id"`
	SenderID     uuid.UUID      `json:"sender_id"`
	RoomID       uuid.UUID      `json:"room_id"`
	ServerTime   int64          `json:"server_time"`
	ReplyToMsgID *uuid.UUID     `json:"reply_to_msg_id"`
	MsgType      string         `json:"msg_type"`
	Payload      jsontext.Value `json:"payload"`
	Ext          jsontext.Value `json:"ext"`
}
