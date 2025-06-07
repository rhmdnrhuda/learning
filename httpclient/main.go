package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// HTTPResponse represents the common response structure from HTTP requests
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Error      error
}

// Create a shared standard HTTP client for better connection reuse
var standardClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		ForceAttemptHTTP2:   false,
	},
}

// StandardGet makes a GET request using the standard net/http package
func StandardGet(ctx context.Context, url string, headers map[string]string, timeout time.Duration) HTTPResponse {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error creating request: %w", err)}
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Use the shared client
	resp, err := standardClient.Do(req)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error making request: %w", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error reading body: %w", err)}
	}

	// Extract headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	return HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    respHeaders,
	}
}

// StandardPost makes a POST request using the standard net/http package
func StandardPost(ctx context.Context, url string, headers map[string]string, body interface{}, timeout time.Duration) HTTPResponse {
	// Marshal body to JSON if it's not nil
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return HTTPResponse{Error: fmt.Errorf("error marshaling request body: %w", err)}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error creating request: %w", err)}
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Set content type if not specified
	if _, exists := headers["Content-Type"]; !exists && body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Use the shared client
	resp, err := standardClient.Do(req)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error making request: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error reading body: %w", err)}
	}

	// Extract headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	return HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    respHeaders,
	}
}

// Create a shared fasthttp client
var fasthttpClient = &fasthttp.Client{
	MaxConnsPerHost:          1000,
	ReadTimeout:              10 * time.Second,
	WriteTimeout:             10 * time.Second,
	NoDefaultUserAgentHeader: true, // Don't add default user-agent
	DisablePathNormalizing:   true,
}

// FastHTTPGet makes a GET request using the fasthttp package
func FastHTTPGet(url string, headers map[string]string, timeout time.Duration) HTTPResponse {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	err := fasthttpClient.DoTimeout(req, resp, timeout)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error making request: %w", err)}
	}

	// Extract headers
	respHeaders := make(map[string]string)
	resp.Header.VisitAll(func(key, value []byte) {
		respHeaders[string(key)] = string(value)
	})

	return HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    respHeaders,
	}
}

// FastHTTPPost makes a POST request using the fasthttp package
func FastHTTPPost(url string, headers map[string]string, body interface{}, timeout time.Duration) HTTPResponse {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set content type if not specified
	if _, exists := headers["Content-Type"]; !exists && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Marshal body to JSON if it's not nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return HTTPResponse{Error: fmt.Errorf("error marshaling request body: %w", err)}
		}
		req.SetBody(bodyBytes)
	}

	err := fasthttpClient.DoTimeout(req, resp, timeout)
	if err != nil {
		return HTTPResponse{Error: fmt.Errorf("error making request: %w", err)}
	}

	// Extract headers
	respHeaders := make(map[string]string)
	resp.Header.VisitAll(func(key, value []byte) {
		respHeaders[string(key)] = string(value)
	})

	return HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    respHeaders,
	}
}
