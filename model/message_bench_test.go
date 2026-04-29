package model

import (
	jsonv2 "encoding/json/v2"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func prepareTestData() (*JsonChatMessage, *ProtoChatMessage) {
	msgID := uuid.New()
	clientMsgID := uuid.New()
	senderID := uuid.New()
	roomID := uuid.New()
	replyToMsgID := uuid.New()
	now := time.Now().UTC()
	payloadBytes := []byte(`{"content":"hello world","mentions":["user1","user2"],"is_rich":true}`)
	extBytes := []byte(`{"source":"ios","version":"1.0.0"}`)

	jsonMsg := &JsonChatMessage{
		MsgID:        msgID,
		ClientMsgID:  clientMsgID,
		SenderID:     senderID,
		RoomID:       roomID,
		ServerTime:   now.UnixMilli(),
		ReplyToMsgID: &replyToMsgID,
		MsgType:      "text",
		Payload:      payloadBytes,
		Ext:          extBytes,
	}

	protoMsg := &ProtoChatMessage{
		MsgId:        msgID[:],
		ClientMsgId:  clientMsgID[:],
		SenderId:     senderID[:],
		RoomId:       roomID[:],
		ServerTime:   now.UnixMilli(),
		ReplyToMsgId: replyToMsgID[:],
		MsgType:      "text",
		Payload:      payloadBytes,
		Ext:          extBytes,
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

	t.Run("手写 Marshal/Unmarshal", func(t *testing.T) {
		jsonData, err := jsonMsg.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var jsonTarget JsonChatMessage
		err = jsonTarget.Unmarshal(jsonData)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(jsonMsg, &jsonTarget) {
			t.Errorf("手写 Marshal/Unmarshal mismatch:\nOriginal: %+v\nUnmarshaled: %+v", jsonMsg, &jsonTarget)
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
	encodeData, err := jsonMsg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JSONv2  Size : %3d bytes", len(jsonData))
	t.Logf("Proto   Size : %3d bytes", len(protoData))
	t.Logf("VTProto Size : %3d bytes", len(vtData))
	t.Logf("手写    Size : %3d bytes", len(encodeData))

	if len(protoData) != len(vtData) {
		t.Errorf("Proto and VTProto sizes mismatch: %d vs %d", len(protoData), len(vtData))
	}
}

func BenchmarkMarshal(b *testing.B) {
	jsonMsg, protoMsg := prepareTestData()
	b.Run("JSONv2 Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, err := jsonv2.Marshal(jsonMsg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Proto Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, err := proto.Marshal(protoMsg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("VTProto Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, err := protoMsg.MarshalVT()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("手写 Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			_, err := jsonMsg.Marshal()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	jsonMsg, protoMsg := prepareTestData()
	jsonData, _ := jsonv2.Marshal(jsonMsg)
	protoData, _ := proto.Marshal(protoMsg)
	vtData, _ := protoMsg.MarshalVT()
	encodeData, _ := jsonMsg.Marshal()

	b.Run("JSONv2 Unmarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			var target JsonChatMessage
			err := jsonv2.Unmarshal(jsonData, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Proto Unmarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			var target ProtoChatMessage
			err := proto.Unmarshal(protoData, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("VTProto Unmarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			var target ProtoChatMessage
			err := target.UnmarshalVT(vtData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("VTProto Unmarshal with Pool", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			target := ProtoChatMessageFromVTPool()
			err := target.UnmarshalVT(vtData)
			if err != nil {
				b.Fatal(err)
			}
			target.ReturnToVTPool()
		}
	})

	b.Run("手写 Unmarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			var target JsonChatMessage
			err := target.Unmarshal(encodeData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("手写 Unmarshal with Pool", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			target := AcquireMessage()
			err := target.Unmarshal(encodeData)
			if err != nil {
				b.Fatal(err)
			}
			ReleaseMessage(target)
		}
	})
}
