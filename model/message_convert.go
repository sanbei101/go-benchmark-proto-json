package model

import (
	"github.com/google/uuid"
)

func (j *JsonChatMessage) ToProto(pb *ProtoChatMessage) *ProtoChatMessage {
	if pb == nil {
		pb = &ProtoChatMessage{}
	}
	pb.MsgId = j.MsgID[:]
	pb.ClientMsgId = j.ClientMsgID[:]
	pb.SenderId = j.SenderID[:]
	if j.ReplyToMsgID != nil {
		pb.ReplyToMsgId = j.ReplyToMsgID[:]
	} else {
		pb.ReplyToMsgId = nil
	}

	pb.MsgType = j.MsgType
	pb.ServerTime = j.ServerTime

	pb.Payload = j.Payload
	pb.Ext = j.Ext
	return pb
}

func (p *ProtoChatMessage) FromProto(j *JsonChatMessage) *JsonChatMessage {
	if j == nil {
		j = &JsonChatMessage{}
	}
	copy(j.MsgID[:], p.MsgId)
	copy(j.ClientMsgID[:], p.ClientMsgId)
	copy(j.SenderID[:], p.SenderId)
	copy(j.RoomID[:], p.RoomId)

	if len(p.ReplyToMsgId) == 16 {
		var replyID uuid.UUID
		copy(replyID[:], p.ReplyToMsgId)
		j.ReplyToMsgID = &replyID
	}

	j.MsgType = p.MsgType
	j.ServerTime = p.ServerTime

	j.Payload = p.Payload
	j.Ext = p.Ext
	return j
}
