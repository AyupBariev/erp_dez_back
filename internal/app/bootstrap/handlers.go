package bootstrap

import (
	. "erp/internal/app/handlers"
	"erp/internal/app/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handlers struct {
	Telegram *TelegramHandler
	Auth     *AuthHandler
	Admin    *AdminHandler
	User     *UserHandler
	Engineer *EngineerHandler
	Order    *OrderHandler
}

func NewHandlers(
	bot *tgbotapi.BotAPI,
	userService *services.UserService,
	engineerService *services.EngineerService,
	orderService *services.OrderService,
	authService *services.AuthService,
) Handlers {
	telegramHandler := NewTelegramHandler(bot, engineerService, orderService)
	authHandler := NewAuthHandler(authService)
	adminHandler := NewAdminHandler(engineerService, telegramHandler)
	userHandler := NewUserHandler(userService)
	engineerHandler := NewEngineerHandler(engineerService)
	orderHandler := NewOrderHandler(orderService, engineerService)

	return Handlers{
		Telegram: telegramHandler,
		Auth:     authHandler,
		Admin:    adminHandler,
		User:     userHandler,
		Engineer: engineerHandler,
		Order:    orderHandler,
	}
}
