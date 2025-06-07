package jsoncompare

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"os"
	"testing"
)

func BenchmarkStandardJSON_SmallFile(b *testing.B) {
	data, err := os.ReadFile("small.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read small.jsoncompare: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []SmallStruct
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkJsonIter_SmallFile(b *testing.B) {
	data, err := os.ReadFile("small.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read small.jsoncompare: %v", err)
	}

	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []SmallStruct
		if err := jsonIter.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkStandardJSON_MediumFile(b *testing.B) {
	data, err := os.ReadFile("medium.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read medium.jsoncompare: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []MediumStruct
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkJsonIter_MediumFile(b *testing.B) {
	data, err := os.ReadFile("medium.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read medium.jsoncompare: %v", err)
	}

	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []MediumStruct
		if err := jsonIter.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkStandardJSON_LargeFile(b *testing.B) {
	data, err := os.ReadFile("large.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read large.jsoncompare: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []LargeStruct
		if err := json.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}

func BenchmarkJsonIter_LargeFile(b *testing.B) {
	data, err := os.ReadFile("large.jsoncompare")
	if err != nil {
		b.Fatalf("Failed to read large.jsoncompare: %v", err)
	}

	var jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []LargeStruct
		if err := jsonIter.Unmarshal(data, &result); err != nil {
			b.Fatalf("Failed to unmarshal: %v", err)
		}
	}
}
