package model

import (
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func prepareTestData() (*JsonChatMessage, *ProtoChatMessage) {
	msgID := uuid.New()
	clientMsgID := uuid.New()
	senderID := uuid.New()
	receiverID := uuid.New()
	replyToMsgID := uuid.New()
	now := time.Now()
	payloadBytes := []byte(`{"content":"hello world","mentions":["user1","user2"],"is_rich":true}`)
	extBytes := []byte(`{"source":"ios","version":"1.0.0"}`)

	jsonMsg := &JsonChatMessage{
		MsgID:        msgID,
		ClientMsgID:  clientMsgID,
		SenderID:     senderID,
		ReceiverID:   receiverID,
		ChatType:     ChatType_CHAT_TYPE_SINGLE,
		ServerTime:   now.UnixMilli(),
		ReplyToMsgID: &replyToMsgID,
		MsgType:      MessageType_MESSAGE_TYPE_TEXT,
		Payload:      payloadBytes,
		Ext:          extBytes,
		CreatedAt:    now,
	}

	protoMsg := &ProtoChatMessage{
		MsgId:        msgID[:],
		ClientMsgId:  clientMsgID[:],
		SenderId:     senderID[:],
		ReceiverId:   receiverID[:],
		ChatType:     ChatType_CHAT_TYPE_SINGLE,
		ServerTime:   now.UnixMilli(),
		ReplyToMsgId: replyToMsgID[:],
		MsgType:      MessageType_MESSAGE_TYPE_TEXT,
		Payload:      payloadBytes,
		Ext:          extBytes,
		CreatedAt:    timestamppb.New(now),
	}

	return jsonMsg, protoMsg
}

// =====================================================================
// Benchmark: Marshal (序列化)
// =====================================================================

func BenchmarkMarshal_JSONv2(b *testing.B) {
	jsonMsg, _ := prepareTestData()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, err := jsonv2.Marshal(jsonMsg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Proto(b *testing.B) {
	_, protoMsg := prepareTestData()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, err := proto.Marshal(protoMsg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_VTProto(b *testing.B) {
	_, protoMsg := prepareTestData()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, err := protoMsg.MarshalVT()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =====================================================================
// Benchmark: Unmarshal (反序列化)
// =====================================================================

func BenchmarkUnmarshal_JSONv2(b *testing.B) {
	jsonMsg, _ := prepareTestData()
	data, err := jsonv2.Marshal(jsonMsg)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var target JsonChatMessage
		err := jsonv2.Unmarshal(data, &target)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Proto(b *testing.B) {
	_, protoMsg := prepareTestData()
	data, err := proto.Marshal(protoMsg)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var target ProtoChatMessage
		err := proto.Unmarshal(data, &target)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_VTProto(b *testing.B) {
	_, protoMsg := prepareTestData()
	data, err := protoMsg.MarshalVT()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var target ProtoChatMessage
		err := target.UnmarshalVT(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
