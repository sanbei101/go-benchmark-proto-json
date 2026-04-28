package model

import (
	jsonv2 "encoding/json/v2"
	"reflect"
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
	now := time.Now().UTC()
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

// / =====================================================================
// Test: Marshal/Unmarshal Correctness (验证序列化和反序列化的正确性)
// =====================================================================
func TestMarshalUnmarshal(t *testing.T) {
	jsonMsg, protoMsg := prepareTestData()
	t.Run("JSONv2 Marshal/Unmarshal", func(t *testing.T) {
		jsonData, err := jsonv2.Marshal(jsonMsg)
		if err != nil {
			t.Fatal(err)
		}
		var jsonTarget JsonChatMessage
		err = jsonv2.Unmarshal(jsonData, &jsonTarget)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(jsonMsg, &jsonTarget) {
			t.Errorf("JSONv2 Marshal/Unmarshal mismatch:\nOriginal: %+v\nUnmarshaled: %+v", jsonMsg, &jsonTarget)
		}
	})

	t.Run("Proto Marshal/Unmarshal", func(t *testing.T) {
		protoData, err := proto.Marshal(protoMsg)
		if err != nil {
			t.Fatal(err)
		}
		var protoTarget ProtoChatMessage
		err = proto.Unmarshal(protoData, &protoTarget)
		if err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(protoMsg, &protoTarget) {
			t.Errorf("Proto Marshal/Unmarshal mismatch:\nOriginal: %+v\nUnmarshaled: %+v", protoMsg, &protoTarget)
		}
	})

	t.Run("VTProto Marshal/Unmarshal", func(t *testing.T) {
		vtData, err := protoMsg.MarshalVT()
		if err != nil {
			t.Fatal(err)
		}
		var vtTarget ProtoChatMessage
		err = vtTarget.UnmarshalVT(vtData)
		if err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(protoMsg, &vtTarget) {
			t.Errorf("VTProto Marshal/Unmarshal mismatch:\nOriginal: %+v\nUnmarshaled: %+v", protoMsg, &vtTarget)
		}
	})
}

// =====================================================================
// Test: Payload Size Comparison (测量序列化后的字节流大小)
// =====================================================================
func TestPayloadSize(t *testing.T) {
	jsonMsg, protoMsg := prepareTestData()
	jsonData, err := jsonv2.Marshal(jsonMsg)
	if err != nil {
		t.Fatal(err)
	}
	protoData, err := proto.Marshal(protoMsg)
	if err != nil {
		t.Fatal(err)
	}
	vtData, err := protoMsg.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JSONv2  Size : %3d bytes", len(jsonData))
	t.Logf("Proto   Size : %3d bytes", len(protoData))
	t.Logf("VTProto Size : %3d bytes", len(vtData))
	if len(protoData) != len(vtData) {
		t.Errorf("Proto and VTProto sizes mismatch: %d vs %d", len(protoData), len(vtData))
	}
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

func BenchmarkUnmarshal_VTPool(b *testing.B) {
	_, protoMsg := prepareTestData()
	data, err := protoMsg.MarshalVT()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		target := ProtoChatMessageFromVTPool()
		err := target.UnmarshalVT(data)
		if err != nil {
			b.Fatal(err)
		}
		target.ReturnToVTPool()
	}
}
