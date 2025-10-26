package services

import (
	"erp/internal/app/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
)

type NotificationService struct {
	Telegram *TelegramService
	Call     *CallService
	Redis    *redis.Client
}

func NewNotificationService(telegram *TelegramService, call *CallService, redis *redis.Client) *NotificationService {
	return &NotificationService{
		Telegram: telegram,
		Call:     call,
		Redis:    redis,
	}
}

func (s *NotificationService) NotifyEngineerNewOrder(order *models.Order, eng *models.Engineer) {
	if eng == nil || eng.TelegramID == 0 {
		log.Printf("Инженер не найден или не имеет Telegram ID")
		return
	}

	text := formatOrderMessage(order)

	// Создаем inline-клавиатуру с кнопками
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Принять", fmt.Sprintf("accept_%d", order.ERPNumber)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("😄 С удовольствием принять", fmt.Sprintf("accept_happy_%d", order.ERPNumber)),
		),
	)

	// Отправляем сообщение через TelegramHandler
	s.Telegram.SendMessageWithKeyboard(eng.TelegramID, text, buttons)

	//telphin
	//n.Call.ScheduleEngineerCall(engineerID, order)

}

// formatOrderMessage — формирует текст для Telegram-сообщения инженеру
func formatOrderMessage(order *models.Order) string {
	// форматируем дату, если нужно
	scheduled := "не указано"
	if !order.ScheduledAt.IsZero() {
		scheduled = order.ScheduledAt.Format("02.01.2006 15:04")
	}

	return fmt.Sprintf(
		"📦 *Новый заказ № %d*\n\n"+
			"📅 Дата и время: %s\n"+
			"🔧 Проблема: %s\n"+
			"👤 Клиент: *%s*\n"+
			"🏠 Адрес: %s\n\n"+
			"Выберите действие ниже:",
		order.ERPNumber,
		scheduled,
		order.Problem,
		order.ClientName,
		order.Address,
	)
}
