package main

import (
	"erp/internal/app/bootstrap"
	"erp/internal/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env")
	}
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}

	logger.SetupLogger()
	log.Println("=== Starting application ===")

	time.Local = loc
	log.Println("Timezone set to Europe/Moscow")

	app := bootstrap.NewApp()

	go app.StartTelegramBot()
	app.StartHTTPServer()
}
