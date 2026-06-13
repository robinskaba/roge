package roblox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type robloxError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type operationResponse struct {
	Done   bool            `json:"done"`
	Path   string          `json:"path"`
	Result json.RawMessage `json:"response,omitempty"`
	Error  *robloxError    `json:"error,omitempty"`
}

func doUpload(req *http.Request) (*operationResponse, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(b))
	}
	var op operationResponse
	if err := json.Unmarshal(b, &op); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &op, nil
}

func pollOperation(apiKey, operationPath string) (*operationResponse, error) {
	operationPath = strings.TrimPrefix(operationPath, "/")
	url := fmt.Sprintf("https://apis.roblox.com/assets/v1/%s", operationPath)

	for attempt := 1; attempt <= 15; attempt++ {
		wait := time.Duration(attempt) * 400 * time.Millisecond
		time.Sleep(wait)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-api-key", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("poll HTTP %d: %s", resp.StatusCode, string(b))
		}
		var op operationResponse
		if err := json.Unmarshal(b, &op); err != nil {
			return nil, fmt.Errorf("parsing poll: %w", err)
		}
		if op.Error != nil {
			return nil, fmt.Errorf("operation error %d: %s", op.Error.Code, op.Error.Message)
		}
		if op.Done {
			return &op, nil
		}
	}
	return nil, fmt.Errorf("operation did not complete after 15 polls")
}

func fetchBinaryFromCdn(location string) ([]byte, error) {
	resp, err := http.Get(location)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve binary from CDN (HTTP %d)", resp.StatusCode)
	}

	binary, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read CDN body: %w", err)
	}

	return binary, nil
}

// authenticated GET request
func Fetch(url string, apiKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch from %s (HTTP %d)", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}
