package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (j *JsonChatMessage) ToProto(pb *ProtoChatMessage) *ProtoChatMessage {
	if pb == nil {
		pb = &ProtoChatMessage{}
	}
	pb.MsgId = j.MsgID[:]
	pb.ClientMsgId = j.ClientMsgID[:]
	pb.SenderId = j.SenderID[:]
	pb.ReceiverId = j.ReceiverID[:]

	if j.ReplyToMsgID != nil {
		pb.ReplyToMsgId = j.ReplyToMsgID[:]
	} else {
		pb.ReplyToMsgId = nil
	}

	pb.ChatType = j.ChatType
	pb.MsgType = j.MsgType
	pb.ServerTime = j.ServerTime

	pb.Payload = j.Payload
	pb.Ext = j.Ext
	if !j.CreatedAt.IsZero() {
		pb.CreatedAt = timestamppb.New(j.CreatedAt)
	} else {
		pb.CreatedAt = nil
	}
	return pb
}

func (p *ProtoChatMessage) FromProto(j *JsonChatMessage) *JsonChatMessage {
	if j == nil {
		j = &JsonChatMessage{}
	}
	copy(j.MsgID[:], p.MsgId)
	copy(j.ClientMsgID[:], p.ClientMsgId)
	copy(j.SenderID[:], p.SenderId)
	copy(j.ReceiverID[:], p.ReceiverId)

	if len(p.ReplyToMsgId) == 16 {
		var replyID uuid.UUID
		copy(replyID[:], p.ReplyToMsgId)
		j.ReplyToMsgID = &replyID
	}

	j.ChatType = p.ChatType
	j.MsgType = p.MsgType
	j.ServerTime = p.ServerTime

	j.Payload = p.Payload
	j.Ext = p.Ext

	if p.CreatedAt != nil {
		j.CreatedAt = p.CreatedAt.AsTime()
	} else {
		j.CreatedAt = time.Time{}
	}
	return j
}
