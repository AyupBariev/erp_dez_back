package main

import (
	"erp/internal/app/bootstrap"
	"erp/internal/pkg/logger"
	"log"
	"time"
)

func main() {
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
