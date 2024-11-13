package utils

import (
	"github.com/google/uuid"
)

func LoadOrCreateUUID() string {
	if currentUUID := loadUUID(); currentUUID != "" {
		return currentUUID
	}
	newUUID := uuid.New().String()
	saveUUID(newUUID)
	return newUUID
}

func loadUUID() string {
	return ""
}

func saveUUID(uuid string) {
}
