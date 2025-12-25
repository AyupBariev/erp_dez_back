package bootstrap

import (
	. "erp/internal/app/handlers"
	"erp/internal/app/services"
	"erp/internal/interfaces/cleanhandler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handlers struct {
	Telegram           *TelegramHandler
	Auth               *AuthHandler
	Admin              *AdminHandler
	User               *UserHandler
	Engineer           *EngineerHandler
	Order              *OrderHandler
	DictHandler        *DictionaryHandler
	Report             *ReportHandler
	RepeatRequest      *RepeatRequestHandler
	Motivation         *MotivationHandler
	EngineerTarget     *EngineerTargetHandler
	EngineerMotivation *EngineerMotivationHandler
	Profit             *ProfitHandler
	AggregatorPayout   *cleanhandler.AggregatorPayoutHandler
}

func NewHandlers(
	bot *tgbotapi.BotAPI,
	userService *services.UserService,
	engineerService *services.EngineerService,
	orderService *services.OrderService,
	authService *services.AuthService,
	dictService *services.DictionaryService,
	reportService *services.ReportService,
	repeatRequestService *services.RepeatRequestService,
	motivationService *services.MotivationService,
	engineerTargetService *services.EngineerTargetService,
	engineerMotivationService *services.EngineerMotivationService,
	aggPayoutHandler *cleanhandler.AggregatorPayoutHandler,
) Handlers {
	telegramHandler := NewTelegramHandler(bot, engineerService, orderService, reportService)
	authHandler := NewAuthHandler(authService)
	adminHandler := NewAdminHandler(engineerService, telegramHandler)
	userHandler := NewUserHandler(userService)
	engineerHandler := NewEngineerHandler(engineerService)
	orderHandler := NewOrderHandler(orderService, engineerService)
	dictHandler := NewDictionaryHandler(dictService)
	reportHandler := NewReportHandler(reportService, orderService, engineerService)
	repeatRequestHandler := NewRepeatRequestHandler(repeatRequestService)

	motivationHandler := NewMotivationHandler(motivationService)
	engineerTargetHandler := NewEngineerTargetHandler(engineerTargetService)
	engineerMotivationHandler := NewEngineerMotivationHandler(engineerMotivationService)
	profitHandle := NewProfitHandler(engineerMotivationService)
	return Handlers{
		Telegram:           telegramHandler,
		Auth:               authHandler,
		Admin:              adminHandler,
		User:               userHandler,
		Engineer:           engineerHandler,
		Order:              orderHandler,
		DictHandler:        dictHandler,
		Report:             reportHandler,
		RepeatRequest:      repeatRequestHandler,
		Motivation:         motivationHandler,
		EngineerTarget:     engineerTargetHandler,
		EngineerMotivation: engineerMotivationHandler,
		Profit:             profitHandle,
		AggregatorPayout:   aggPayoutHandler,
	}
}
