package model

import (
	"encoding/binary"
	"encoding/json/jsontext"
	"errors"
	sync "sync"

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

var MessagePool = sync.Pool{
	New: func() any {
		return &JsonChatMessage{}
	},
}

func AcquireMessage() *JsonChatMessage {
	return MessagePool.Get().(*JsonChatMessage)
}

func ReleaseMessage(m *JsonChatMessage) {
	m.Reset()
	MessagePool.Put(m)
}

func (m *JsonChatMessage) Reset() {
	m.MsgID = uuid.UUID{}
	m.ClientMsgID = uuid.UUID{}
	m.SenderID = uuid.UUID{}
	m.RoomID = uuid.UUID{}
	m.ServerTime = 0
	m.ReplyToMsgID = nil
	m.MsgType = ""

	// 核心：复用底层数组，长度置 0，容量保留
	if m.Payload != nil {
		m.Payload = m.Payload[:0]
	}
	if m.Ext != nil {
		m.Ext = m.Ext[:0]
	}
}

// Marshal 将结构体编码为紧凑的 []byte
func (m *JsonChatMessage) Marshal() ([]byte, error) {
	// 1. 预计算总容量，避免切片在 append 时频繁扩容（极致性能）
	// MsgID(16) + ClientMsgID(16) + SenderID(16) + RoomID(16) + ServerTime(8)
	// + ReplyToMsgID Flag(1) = 73 字节的固定长度
	size := 16*4 + 8 + 1

	if m.ReplyToMsgID != nil {
		size += 16
	}

	// 可变长度计算: MsgType(2字节长度 + 内容), Payload(4字节长度 + 内容), Ext(4字节长度 + 内容)
	size += 2 + len(m.MsgType)
	size += 4 + len(m.Payload)
	size += 4 + len(m.Ext)

	// 2. 分配精确大小的内存
	buf := make([]byte, size)
	offset := 0

	// 3. 依次写入数据
	copy(buf[offset:], m.MsgID[:])
	offset += 16
	copy(buf[offset:], m.ClientMsgID[:])
	offset += 16
	copy(buf[offset:], m.SenderID[:])
	offset += 16
	copy(buf[offset:], m.RoomID[:])
	offset += 16

	binary.BigEndian.PutUint64(buf[offset:], uint64(m.ServerTime))
	offset += 8

	// 处理指针
	if m.ReplyToMsgID != nil {
		buf[offset] = 1
		offset += 1
		copy(buf[offset:], (*m.ReplyToMsgID)[:])
		offset += 16
	} else {
		buf[offset] = 0
		offset += 1
	}

	// 处理字符串
	binary.BigEndian.PutUint16(buf[offset:], uint16(len(m.MsgType)))
	offset += 2
	copy(buf[offset:], m.MsgType)
	offset += len(m.MsgType)

	// 处理 jsontext.Value ([]byte)
	binary.BigEndian.PutUint32(buf[offset:], uint32(len(m.Payload)))
	offset += 4
	copy(buf[offset:], m.Payload)
	offset += len(m.Payload)

	binary.BigEndian.PutUint32(buf[offset:], uint32(len(m.Ext)))
	offset += 4
	copy(buf[offset:], m.Ext)
	offset += len(m.Ext)

	return buf, nil
}

// Unmarshal 从 []byte 还原结构体 (配合 Pool 使用)
func (m *JsonChatMessage) Unmarshal(data []byte) error {
	if len(data) < 83 {
		return errors.New("data too short")
	}
	offset := 0

	// 1. 固定长度读取 (同前)
	copy(m.MsgID[:], data[offset:offset+16])
	offset += 16
	copy(m.ClientMsgID[:], data[offset:offset+16])
	offset += 16
	copy(m.SenderID[:], data[offset:offset+16])
	offset += 16
	copy(m.RoomID[:], data[offset:offset+16])
	offset += 16

	m.ServerTime = int64(binary.BigEndian.Uint64(data[offset : offset+8]))
	offset += 8

	// 2. 指针读取
	if data[offset] == 1 {
		offset += 1
		if len(data) < offset+16 {
			return errors.New("data too short for ReplyToMsgID")
		}
		id := uuid.UUID{}
		copy(id[:], data[offset:offset+16])
		offset += 16
		m.ReplyToMsgID = &id
	} else {
		offset += 1
		m.ReplyToMsgID = nil
	}

	// 3. MsgType
	msgTypeLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	offset += 2
	m.MsgType = string(data[offset : offset+msgTypeLen])
	offset += msgTypeLen

	// 4. 读取 Payload (复用容量)
	payloadLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if cap(m.Payload) >= payloadLen {
		m.Payload = m.Payload[:payloadLen] // 容量够，直接切片扩充长度
	} else {
		m.Payload = make(jsontext.Value, payloadLen) // 容量不够才 make
	}
	copy(m.Payload, data[offset:offset+payloadLen])
	offset += payloadLen

	// 5. 读取 Ext (复用容量)
	extLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
	offset += 4
	if cap(m.Ext) >= extLen {
		m.Ext = m.Ext[:extLen]
	} else {
		m.Ext = make(jsontext.Value, extLen)
	}
	copy(m.Ext, data[offset:offset+extLen])
	offset += extLen

	return nil
}
