package bootstrap

import (
	"log"
	"os"
)

func (app *App) StartHTTPServer() {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on :%s", port)
	if err := app.HTTPServer.Run(":" + port); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}
