package services

import (
	"erp/internal/app/models"
	"erp/internal/pkg/logger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"strings"
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
	if id := eng.GetTelegramID(); id == nil {
		logger.LogInfo(fmt.Sprintf("Инженер %s не имеет Telegram ID — уведомление не отправлено", eng.Username))
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
	id := eng.GetTelegramID()
	s.Telegram.SendMessageWithKeyboard(*id, text, buttons)

	//telphin
	//n.Call.ScheduleEngineerCall(engineerID, order)

}

func (s *NotificationService) NotifyEngineerOrderUnassigned(order *models.Order, engineer *models.Engineer) {
	if id := engineer.GetTelegramID(); id == nil {
		logger.LogInfo(fmt.Sprintf("Инженер %s не имеет Telegram ID — уведомление не отправлено", engineer.Username))
		return
	}

	message := fmt.Sprintf(
		"❌ Заказ №%d снят с вас\n\n"+
			"📋 Детали заказа:\n"+
			"• Клиент: %s\n"+
			"• Адрес: %s\n"+
			"• Проблема: %s\n"+
			"• Время: %s\n\n"+
			"Заказ возвращен в пул нераспределенных заказов.",
		order.ERPNumber,
		order.ClientName,
		order.Address,
		order.Problem.Name,
		order.ScheduledAt.Format("02.01.2006 15:04"),
	)

	id := engineer.GetTelegramID()
	s.Telegram.sendMessage(*id, message)
}

// formatOrderMessage — формирует текст для Telegram-сообщения инженеру
func formatOrderMessage(order *models.Order) string {
	// форматируем дату, если нужно
	scheduled := "не указано"
	if !order.ScheduledAt.IsZero() {
		scheduled = order.ScheduledAt.Format("02.01.2006 15:04")
	}

	phones := order.Phones
	var phoneDisplay string
	if len(phones) == 0 {
		phoneDisplay = "—"
	} else {
		phoneDisplay = strings.Join(phones, ", ")
	}

	price := order.Price
	if price == "" {
		price = "—"
	}

	return fmt.Sprintf(
		"📦 *Новый заказ № %d*\n\n"+
			"📅 Дата и время: %s\n"+
			"🔧 Проблема: %s\n\n"+
			"🏠 Адрес: %s\n"+
			"👤 Клиент: *%s*\n"+
			"📞 Телефон: %s\n\n"+
			"💰Сумма: %s\n\n"+
			"Выберите действие ниже:",
		order.ERPNumber,
		scheduled,
		order.Problem.Name,
		order.Address,
		order.ClientName,
		phoneDisplay,
		price,
	)
}
