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

const MaxRetries = 3

func MakePOSTRequest(url string, data interface{}) ([]byte, int, error) {
	return makeRequestWithRetry("POST", url, data)
}

func MakeGETRequest(url string) ([]byte, int, error) {
	return makeRequestWithRetry("GET", url, nil)
}

func makeRequestWithRetry(method, url string, data interface{}) ([]byte, int, error) {
	var err error
	var resp *http.Response

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		var req *http.Request
		if method == "POST" && data != nil {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, 0, err
			}
			req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return nil, 0, err
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			return data, resp.StatusCode, err
		}

		if attempt < MaxRetries {
			logger.Warn("Request failed, retrying...", "Attempt:", attempt, "Error:", err)
			time.Sleep(time.Second * time.Duration(attempt*2)) // Exponential backoff
		} else {
			logger.Error("Request failed after max retries:", err)
		}
	}

	return nil, 0, errors.New("failed to make request after retries")
}
