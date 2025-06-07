package jsoncompare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func init() {
	// Generate sample files before running benchmarks
	generateSampleFiles()
}

type SmallStruct struct {
	ID   int    `jsoncompare:"id"`
	Name string `jsoncompare:"name"`
}

type MediumStruct struct {
	ID        int      `jsoncompare:"id"`
	Name      string   `jsoncompare:"name"`
	Email     string   `jsoncompare:"email"`
	Age       int      `jsoncompare:"age"`
	Active    bool     `jsoncompare:"active"`
	CreatedAt string   `jsoncompare:"created_at"`
	Tags      []string `jsoncompare:"tags"`
}

type LargeStruct struct {
	ID        int                    `jsoncompare:"id"`
	Name      string                 `jsoncompare:"name"`
	Email     string                 `jsoncompare:"email"`
	Age       int                    `jsoncompare:"age"`
	Active    bool                   `jsoncompare:"active"`
	CreatedAt string                 `jsoncompare:"created_at"`
	Tags      []string               `jsoncompare:"tags"`
	Address   Address                `jsoncompare:"address"`
	Friends   []Friend               `jsoncompare:"friends"`
	Settings  map[string]interface{} `jsoncompare:"settings"`
}

type Address struct {
	Street  string `jsoncompare:"street"`
	City    string `jsoncompare:"city"`
	State   string `jsoncompare:"state"`
	Zip     string `jsoncompare:"zip"`
	Country string `jsoncompare:"country"`
}

type Friend struct {
	ID    int    `jsoncompare:"id"`
	Name  string `jsoncompare:"name"`
	Email string `jsoncompare:"email"`
}

// Benchmark function
func benchmark(name string, data []byte, v interface{}, fn func([]byte, interface{}) error) {
	fmt.Printf("Benchmarking %s...\n", name)
	iterations := 1000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		err := fn(data, v)
		if err != nil {
			fmt.Printf("Error unmarshaling: %v\n", err)
			return
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("%s took %s (avg: %s per op)\n", name, elapsed, elapsed/time.Duration(iterations))
}

// Generate sample files with different sizes
func generateSampleFiles() {
	// Small file
	small := []SmallStruct{}
	for i := 0; i < 10; i++ {
		small = append(small, SmallStruct{ID: i, Name: fmt.Sprintf("Name %d", i)})
	}
	smallData, _ := json.Marshal(small)
	ioutil.WriteFile("small.jsoncompare", smallData, 0644)

	// Medium file
	medium := []MediumStruct{}
	for i := 0; i < 100; i++ {
		medium = append(medium, MediumStruct{
			ID:        i,
			Name:      fmt.Sprintf("Name %d", i),
			Email:     fmt.Sprintf("email%d@example.com", i),
			Age:       20 + (i % 50),
			Active:    i%2 == 0,
			CreatedAt: time.Now().Format(time.RFC3339),
			Tags:      []string{"tag1", "tag2", "tag3"},
		})
	}
	mediumData, _ := json.Marshal(medium)
	ioutil.WriteFile("medium.jsoncompare", mediumData, 0644)

	// Large file
	large := []LargeStruct{}
	for i := 0; i < 1000; i++ {
		friends := []Friend{}
		for j := 0; j < 5; j++ {
			friends = append(friends, Friend{
				ID:    j,
				Name:  fmt.Sprintf("Friend %d", j),
				Email: fmt.Sprintf("friend%d@example.com", j),
			})
		}

		settings := map[string]interface{}{
			"notification": true,
			"theme":        "dark",
			"timezone":     "UTC",
			"language":     "en",
			"preferences": map[string]interface{}{
				"autoSave":  true,
				"fontSize":  12,
				"fontColor": "#333333",
			},
		}

		large = append(large, LargeStruct{
			ID:        i,
			Name:      fmt.Sprintf("Name %d", i),
			Email:     fmt.Sprintf("email%d@example.com", i),
			Age:       20 + (i % 50),
			Active:    i%2 == 0,
			CreatedAt: time.Now().Format(time.RFC3339),
			Tags:      []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
			Address: Address{
				Street:  fmt.Sprintf("%d Main St", i),
				City:    "Anytown",
				State:   "ST",
				Zip:     "12345",
				Country: "Country",
			},
			Friends:  friends,
			Settings: settings,
		})
	}
	largeData, _ := json.Marshal(large)
	ioutil.WriteFile("large.jsoncompare", largeData, 0644)
}

func benchmarkFile(filename string) {
	fmt.Printf("\n=== Benchmarking %s ===\n", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	switch filename {
	case "small.jsoncompare":
		var standardResult []SmallStruct
		var jsoniterResult []SmallStruct
		benchmark("encoding/json", data, &standardResult, StandardJsonUnmarshal)
		benchmark("jsoniter", data, &jsoniterResult, JsonIterUnmarshal)

	case "medium.jsoncompare":
		var standardResult []MediumStruct
		var jsoniterResult []MediumStruct
		benchmark("encoding/json", data, &standardResult, StandardJsonUnmarshal)
		benchmark("jsoniter", data, &jsoniterResult, JsonIterUnmarshal)

	case "large.jsoncompare":
		var standardResult []LargeStruct
		var jsoniterResult []LargeStruct
		benchmark("encoding/json", data, &standardResult, StandardJsonUnmarshal)
		benchmark("jsoniter", data, &jsoniterResult, JsonIterUnmarshal)
	}
}

func main() {
	fmt.Println("Generating sample JSON files...")
	generateSampleFiles()

	benchmarkFile("small.jsoncompare")
	benchmarkFile("medium.jsoncompare")
	benchmarkFile("large.jsoncompare")
}
