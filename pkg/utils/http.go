package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"kiosk-client/pkg/logger"
	"net/http"
	"time"
)

const MaxRetries = 3

func MakePOSTRequest(url string, data interface{}) ([]byte, int, error) {
	return makeRequestWithRetry("POST", url, data)
}

func MakeGETRequest(url string) ([]byte, int, error) {
	return makeRequestWithRetry("GET", url, nil)
}

func MakePUTRequest(url string, data interface{}) ([]byte, int, error) {
	return makeRequestWithRetry("PUT", url, data)
}

func makeRequestWithRetry(method, url string, data interface{}) ([]byte, int, error) {
	var err error
	var resp *http.Response

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		var req *http.Request
		if data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, 0, err
			}
			req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, 0, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, url, nil)
			if err != nil {
				return nil, 0, err
			}
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err = client.Do(req)

		if err != nil {
			if attempt < MaxRetries {
				logger.Warn("Request failed, retrying...", "Attempt:", attempt, "Error:", err)
				time.Sleep(time.Second * time.Duration(attempt*2)) // Exponential backoff
				continue
			} else {
				logger.Error("Request failed after max retries:", err)
				return nil, 0, errors.New("failed to make request after retries")
			}
		}

		defer resp.Body.Close()

		// Read and return the response if successful
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			return body, resp.StatusCode, err
		} else {
			if attempt < MaxRetries {
				logger.Warn("Non-200 status, retrying...", "Attempt:", attempt, "StatusCode:", resp.StatusCode)
				time.Sleep(time.Second * time.Duration(attempt*2)) // Exponential backoff
			} else {
				logger.Error("Request failed after max retries with status code:", resp.StatusCode)
				return nil, resp.StatusCode, errors.New("failed to make request after retries")
			}
		}
	}

	return nil, 0, errors.New("exceeded retry limit")
}
