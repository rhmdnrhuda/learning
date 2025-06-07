package httpclient

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkResponseSizeComparison(b *testing.B) {
	// Small payload tests
	b.Run("Small-Gzip", func(b *testing.B) {
		server := setupCompressedTestServer("small")
		defer server.Close()

		headers := map[string]string{
			"User-Agent":      "Benchmark-Client",
			"Accept-Encoding": "gzip",
		}
		timeout := 5 * time.Second

		var totalBytes int64
		var wireBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))

			// Try to get wire bytes from Content-Length if available
			if s, ok := resp.Headers["Content-Length"]; ok {
				var cl int
				fmt.Sscanf(s, "%d", &cl)
				wireBytes += int64(cl)
			}
		}

		// Use consistent metric names
		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")

		// If we have wire bytes, report compression metrics
		if wireBytes > 0 {
			b.ReportMetric(float64(wireBytes)/float64(b.N), "wire_bytes/op")
			compressionRatio := float64(totalBytes) / float64(wireBytes)
			b.ReportMetric(compressionRatio, "compression_ratio")
		}
	})

	b.Run("Small-NoGzip", func(b *testing.B) {
		server := setupUncompressedTestServer("small")
		defer server.Close()

		headers := map[string]string{
			"User-Agent": "Benchmark-Client",
		}
		timeout := 5 * time.Second

		var totalBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))
		}

		// Use consistent metric names (same as gzip version)
		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")
		// For uncompressed, wire bytes = decoded bytes
		b.ReportMetric(float64(totalBytes)/float64(b.N), "wire_bytes/op")
		// Compression ratio is 1.0 for uncompressed
		b.ReportMetric(1.0, "compression_ratio")
	})

	// Medium payload tests
	b.Run("Medium-Gzip", func(b *testing.B) {
		server := setupCompressedTestServer("medium")
		defer server.Close()

		headers := map[string]string{
			"User-Agent":      "Benchmark-Client",
			"Accept-Encoding": "gzip",
		}
		timeout := 5 * time.Second

		var totalBytes int64
		var wireBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))

			if s, ok := resp.Headers["Content-Length"]; ok {
				var cl int
				fmt.Sscanf(s, "%d", &cl)
				wireBytes += int64(cl)
			}
		}

		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")

		if wireBytes > 0 {
			b.ReportMetric(float64(wireBytes)/float64(b.N), "wire_bytes/op")
			compressionRatio := float64(totalBytes) / float64(wireBytes)
			b.ReportMetric(compressionRatio, "compression_ratio")
		}
	})

	b.Run("Medium-NoGzip", func(b *testing.B) {
		server := setupUncompressedTestServer("medium")
		defer server.Close()

		headers := map[string]string{
			"User-Agent": "Benchmark-Client",
		}
		timeout := 5 * time.Second

		var totalBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))
		}

		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")
		b.ReportMetric(float64(totalBytes)/float64(b.N), "wire_bytes/op")
		b.ReportMetric(1.0, "compression_ratio")
	})

	// Large payload tests
	b.Run("Large-Gzip", func(b *testing.B) {
		server := setupCompressedTestServer("large")
		defer server.Close()

		headers := map[string]string{
			"User-Agent":      "Benchmark-Client",
			"Accept-Encoding": "gzip",
		}
		timeout := 5 * time.Second

		var totalBytes int64
		var wireBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))

			if s, ok := resp.Headers["Content-Length"]; ok {
				var cl int
				fmt.Sscanf(s, "%d", &cl)
				wireBytes += int64(cl)
			}
		}

		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")

		if wireBytes > 0 {
			b.ReportMetric(float64(wireBytes)/float64(b.N), "wire_bytes/op")
			compressionRatio := float64(totalBytes) / float64(wireBytes)
			b.ReportMetric(compressionRatio, "compression_ratio")
		}
	})

	b.Run("Large-NoGzip", func(b *testing.B) {
		server := setupUncompressedTestServer("large")
		defer server.Close()

		headers := map[string]string{
			"User-Agent": "Benchmark-Client",
		}
		timeout := 5 * time.Second

		var totalBytes int64

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
			totalBytes += int64(len(resp.Body))
		}

		b.ReportMetric(float64(totalBytes)/float64(b.N), "decoded_bytes/op")
		b.ReportMetric(float64(totalBytes)/float64(b.N), "wire_bytes/op")
		b.ReportMetric(1.0, "compression_ratio")
	})
}

// setupCompressedTestServer creates a test server with gzip compression
func setupCompressedTestServer(jsonFile string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")

		// Create a gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Write the JSON content
		var jsonContent []byte
		switch jsonFile {
		case "small":
			jsonContent = []byte(`{"message":"success","count":1}`)
		case "medium":
			// ~1KB of JSON data
			jsonContent = generateMediumJSON()
		case "large":
			// ~10KB of JSON data
			jsonContent = generateLargeJSON()
		}

		w.WriteHeader(http.StatusOK)
		gz.Write(jsonContent)
	})

	server := httptest.NewServer(handler)
	server.EnableHTTP2 = false
	return server
}

// setupUncompressedTestServer creates a test server without compression
func setupUncompressedTestServer(jsonFile string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON content
		var jsonContent []byte
		switch jsonFile {
		case "small":
			jsonContent = []byte(`{"message":"success","count":1}`)
		case "medium":
			// ~1KB of JSON data
			jsonContent = generateMediumJSON()
		case "large":
			// ~10KB of JSON data
			jsonContent = generateLargeJSON()
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonContent)
	})

	server := httptest.NewServer(handler)
	server.EnableHTTP2 = false
	return server
}

// Helper functions to generate different payload sizes
func generateMediumJSON() []byte {
	items := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		items[i] = map[string]interface{}{
			"id":          i,
			"name":        fmt.Sprintf("Item %d", i),
			"description": fmt.Sprintf("This is item %d with some additional text to increase payload size", i),
			"tags":        []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
			"metadata": map[string]interface{}{
				"createdAt": time.Now().Format(time.RFC3339),
				"updatedAt": time.Now().Format(time.RFC3339),
				"active":    true,
				"priority":  i % 5,
			},
		}
	}

	data := map[string]interface{}{
		"items":  items,
		"count":  len(items),
		"status": "success",
	}

	jsonBytes, _ := json.Marshal(data)
	return jsonBytes
}

func generateLargeJSON() []byte {
	items := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		items[i] = map[string]interface{}{
			"id":          i,
			"uuid":        fmt.Sprintf("uuid-%d-%d", i, time.Now().UnixNano()),
			"name":        fmt.Sprintf("Item %d", i),
			"description": fmt.Sprintf("This is item %d with a much longer description to increase payload size substantially. Including additional text with repeated information to make it even larger.", i),
			"tags":        []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8"},
			"metadata": map[string]interface{}{
				"createdAt":      time.Now().Format(time.RFC3339),
				"updatedAt":      time.Now().Format(time.RFC3339),
				"active":         true,
				"priority":       i % 5,
				"category":       fmt.Sprintf("Category %d", i%10),
				"subcategory":    fmt.Sprintf("Subcategory %d", i%20),
				"views":          i * 100,
				"favoriteCount":  i % 50,
				"commentCount":   i % 25,
				"lastModifiedBy": fmt.Sprintf("user-%d", i%15),
				"regions":        []string{"us-east", "us-west", "eu-central", "ap-south"},
			},
			"details": map[string]interface{}{
				"manufacturer": fmt.Sprintf("Company %d", i%10),
				"origin":       fmt.Sprintf("Country %d", i%30),
				"year":         2020 + (i % 5),
				"dimensions": map[string]interface{}{
					"width":  10.5 + float64(i%10),
					"height": 20.5 + float64(i%15),
					"depth":  5.5 + float64(i%8),
					"weight": 2.5 + float64(i%10),
				},
			},
		}
	}

	data := map[string]interface{}{
		"items":       items,
		"count":       len(items),
		"status":      "success",
		"totalPages":  10,
		"currentPage": 1,
		"pageSize":    100,
		"metadata": map[string]interface{}{
			"apiVersion": "1.0.0",
			"timestamp":  time.Now().Format(time.RFC3339),
		},
	}

	jsonBytes, _ := json.Marshal(data)
	return jsonBytes
}

func BenchmarkFastHTTPGzipSmall(b *testing.B) {
	server := setupCompressedTestServer("small")
	defer server.Close()

	headers := map[string]string{
		"User-Agent":      "Benchmark-Client",
		"Accept-Encoding": "gzip",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPNoGzipSmall(b *testing.B) {
	server := setupUncompressedTestServer("small")
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPGzipMedium(b *testing.B) {
	server := setupCompressedTestServer("medium")
	defer server.Close()

	headers := map[string]string{
		"User-Agent":      "Benchmark-Client",
		"Accept-Encoding": "gzip",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPNoGzipMedium(b *testing.B) {
	server := setupUncompressedTestServer("medium")
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPGzipLarge(b *testing.B) {
	server := setupCompressedTestServer("large")
	defer server.Close()

	headers := map[string]string{
		"User-Agent":      "Benchmark-Client",
		"Accept-Encoding": "gzip",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPNoGzipLarge(b *testing.B) {
	server := setupUncompressedTestServer("large")
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPGet(server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

// Parallel versions
func BenchmarkFastHTTPGzipLarge_Parallel(b *testing.B) {
	server := setupCompressedTestServer("large")
	defer server.Close()

	headers := map[string]string{
		"User-Agent":      "Benchmark-Client",
		"Accept-Encoding": "gzip",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
		}
	})
}

func BenchmarkFastHTTPNoGzipLarge_Parallel(b *testing.B) {
	server := setupUncompressedTestServer("large")
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := FastHTTPGet(server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
		}
	})
}
