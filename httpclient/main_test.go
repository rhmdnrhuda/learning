package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setupTestServer creates a mock server for benchmarking
func setupTestServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})
	server := httptest.NewServer(handler)

	// Force HTTP/1.1 to ensure compatibility with fasthttp
	server.EnableHTTP2 = false
	return server
}

// Create a single test server for all parallel benchmarks
var parallelServer *httptest.Server

func init() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})
	parallelServer = httptest.NewServer(handler)
	parallelServer.EnableHTTP2 = false
}

func BenchmarkStandardGet(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	ctx := context.Background()
	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := StandardGet(ctx, server.URL, headers, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkStandardPost(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	ctx := context.Background()
	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	body := map[string]interface{}{
		"test": "data",
		"num":  123,
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := StandardPost(ctx, server.URL, headers, body, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

func BenchmarkFastHTTPGet(b *testing.B) {
	server := setupTestServer()
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

func BenchmarkFastHTTPPost(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	body := map[string]interface{}{
		"test": "data",
		"num":  123,
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp := FastHTTPPost(server.URL, headers, body, timeout)
		if resp.Error != nil {
			b.Fatal(resp.Error)
		}
	}
}

// Parallel benchmarks to simulate concurrent requests

func BenchmarkStandardGet_Parallel(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	ctx := context.Background()
	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := StandardGet(ctx, server.URL, headers, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
		}
	})
}

func BenchmarkStandardPost_Parallel(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	ctx := context.Background()
	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	body := map[string]interface{}{
		"test": "data",
		"num":  123,
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := StandardPost(ctx, server.URL, headers, body, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
		}
	})
}

func BenchmarkFastHTTPGet_Parallel(b *testing.B) {
	server := setupTestServer()
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

func BenchmarkFastHTTPPost_Parallel(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	headers := map[string]string{
		"User-Agent": "Benchmark-Client",
	}
	body := map[string]interface{}{
		"test": "data",
		"num":  123,
	}
	timeout := 5 * time.Second

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp := FastHTTPPost(server.URL, headers, body, timeout)
			if resp.Error != nil {
				b.Fatal(resp.Error)
			}
		}
	})
}
