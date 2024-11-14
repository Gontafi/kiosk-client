package utils

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

const uuidFilePath = "device_uuid.txt"

func LoadOrCreateUUID() string {
	if currentUUID := loadUUID(); currentUUID != "" {
		return currentUUID
	}
	newUUID := uuid.New().String()
	saveUUID(newUUID)
	return newUUID
}

func loadUUID() string {
	data, err := os.ReadFile(uuidFilePath)
	if err != nil {
		return ""
	}
	return string(data)
}

func saveUUID(uuid string) {
	err := os.WriteFile(uuidFilePath, []byte(uuid), 0644)
	if err != nil {
		fmt.Printf("Failed to save UUID to file: %v\n", err)
	}
}
