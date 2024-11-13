package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"kiosk-client/pkg/logger"
	"net/http"
	"time"
)

// MaxRetries is the number of times to retry a failed HTTP request.
const MaxRetries = 3

// MakePOSTRequest sends a POST request to the specified URL with JSON payload and retries on failure.
func MakePOSTRequest(url string, data interface{}) ([]byte, error) {
	return makeRequestWithRetry("POST", url, data)
}

// MakeGETRequest sends a GET request to the specified URL and retries on failure.
func MakeGETRequest(url string) ([]byte, error) {
	return makeRequestWithRetry("GET", url, nil)
}

// makeRequestWithRetry sends an HTTP request with retry logic for resilience.
func makeRequestWithRetry(method, url string, data interface{}) ([]byte, error) {
	var err error
	var resp *http.Response

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Create the HTTP request with JSON body if data is provided
		var req *http.Request
		if method == "POST" && data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return nil, err
		}

		// Send the request and check for response
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			return ioutil.ReadAll(resp.Body)
		}

		// Log error and retry if not last attempt
		if attempt < MaxRetries {
			logger.Warn("Request failed, retrying...", "Attempt:", attempt, "Error:", err)
			time.Sleep(time.Second * time.Duration(attempt*2)) // Exponential backoff
		} else {
			logger.Error("Request failed after max retries:", err)
		}
	}

	return nil, errors.New("failed to make request after retries")
}

// ParseJSONResponse parses the JSON response into the provided struct.
func ParseJSONResponse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
