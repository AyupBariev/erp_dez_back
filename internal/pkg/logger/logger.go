package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func SetupLogger() {
	// Создание каталога logs, если не существует
	_ = os.MkdirAll("logs", os.ModePerm)

	// Форматируем имя файла по текущей дате
	currentDate := time.Now().Format("2006-01-02")
	logFile := fmt.Sprintf("logs/app-%s.log", currentDate)

	// Открываем файл на дозапись
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))
}

func LogError(message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	}
}

func LogInfo(message string) {
	log.Printf("[INFO] %s", message)
}
