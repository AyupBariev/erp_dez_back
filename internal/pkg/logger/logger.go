package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

func SetupLogger() {
	_ = os.MkdirAll("logs", os.ModePerm)
	currentDate := time.Now().Format("2006-01-02")
	logFile := fmt.Sprintf("logs/app-%s.log", currentDate)

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LogError(message string, err error) {
	if err != nil {
		// Получаем место вызова
		_, file, line, _ := runtime.Caller(1)
		log.Printf("[ERROR] %s:%d %s: %v", file, line, message, err)
	}
}

func LogInfo(message string) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("[INFO] %s:%d %s", file, line, message)
}
