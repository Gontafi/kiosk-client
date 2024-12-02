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

const MaxRetries = 2

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
    for attempt := 1; attempt <= MaxRetries; attempt++ {
        var req *http.Request
        var err error

        if data != nil {
            jsonData, err := json.Marshal(data)
            if err != nil {
                return nil, 0, err
            }
            req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
        } else {
            req, err = http.NewRequest(method, url, nil)
        }

        if err != nil {
            return nil, 0, err
        }

        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Timeout: 3 * time.Second}

        resp, err := client.Do(req)
        if err != nil {
            if attempt < MaxRetries {
                logger.Warn("Request failed, retrying...", "Attempt:", attempt, "Error:", err)
                time.Sleep(time.Second * time.Duration(attempt))
                continue
            }
            logger.Error("Request failed after max retries:", err)
            return nil, 0, errors.New("failed to make request after retries")
        }

        body, err := io.ReadAll(resp.Body)
        resp.Body.Close() // Explicitly close the body
        if err != nil {
            return nil, resp.StatusCode, err
        }

        if resp.StatusCode == http.StatusOK {
            return body, resp.StatusCode, nil
        }

        if attempt < MaxRetries {
            logger.Warn("Non-200 status, retrying...", "Attempt:", attempt, "StatusCode:", resp.StatusCode)
            time.Sleep(time.Second * time.Duration(attempt))
        } else {
            logger.Error("Request failed after max retries with status code:", resp.StatusCode)
            return nil, resp.StatusCode, errors.New("failed to make request after retries")
        }
    }

    return nil, 0, errors.New("exceeded retry limit")
}
