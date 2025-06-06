package main

import (
	"encoding/json"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	// Generate sample files before running benchmarks
	generateSampleFiles()
}

func BenchmarkStandardJSON_SmallFile(b *testing.B) {
	data, err := os.ReadFile("small.json")
	if err != nil {
		b.Fatalf("Failed to read small.json: %v", err)
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
	data, err := os.ReadFile("small.json")
	if err != nil {
		b.Fatalf("Failed to read small.json: %v", err)
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
	data, err := os.ReadFile("medium.json")
	if err != nil {
		b.Fatalf("Failed to read medium.json: %v", err)
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
	data, err := os.ReadFile("medium.json")
	if err != nil {
		b.Fatalf("Failed to read medium.json: %v", err)
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
	data, err := os.ReadFile("large.json")
	if err != nil {
		b.Fatalf("Failed to read large.json: %v", err)
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
	data, err := os.ReadFile("large.json")
	if err != nil {
		b.Fatalf("Failed to read large.json: %v", err)
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
