package model

import (
	reflect "reflect"
	"testing"
)

func TestConversion(t *testing.T) {
	jsonMsg, _ := prepareTestData()

	pb := jsonMsg.ToProto(nil)

	recoveredJson := pb.FromProto(nil)

	if !reflect.DeepEqual(jsonMsg, recoveredJson) {
		t.Errorf("互相转换后数据丢失或不匹配!\n原数据: %+v\n新数据: %+v", jsonMsg, recoveredJson)
	} else {
		t.Log("互相转换测试通过,数据完美一致!")
	}
}

// --- Benchmark: Json 转 Proto ---
func BenchmarkConvert_ToProto_New(b *testing.B) {
	jsonMsg, _ := prepareTestData()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = jsonMsg.ToProto(nil)
	}
}

func BenchmarkConvert_ToProto_Reuse(b *testing.B) {
	jsonMsg, _ := prepareTestData()
	reusablePb := &ProtoChatMessage{}
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = jsonMsg.ToProto(reusablePb)
	}
}

func BenchmarkConvert_FromProto_New(b *testing.B) {
	_, protoMsg := prepareTestData()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = protoMsg.FromProto(nil)
	}
}

func BenchmarkConvert_FromProto_Reuse(b *testing.B) {
	_, protoMsg := prepareTestData()
	reusableJson := &JsonChatMessage{}
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = protoMsg.FromProto(reusableJson)
	}
}
