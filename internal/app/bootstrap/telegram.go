package bootstrap

import "log"

func (app *App) StartTelegramBot() {
	log.Println("Telegram bot started")
	go app.Handlers.Telegram.HandleUpdates()
}
