package logger

import (
	"log"
)

func Info(args ...interface{}) {
	log.Println("[INFO]", args)
}

func Error(args ...interface{}) {
	log.Println("[ERROR]", args)
}

func Warn(args ...interface{}) {
	log.Println("[WARN]", args)
}
