package handlers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *TelegramHandler) showOrdersMenu(chatID int64) {
	text := "*Мои заказы\\. *\nСписок актуальных заказов:"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Сегодня", "orders_today"),
			tgbotapi.NewInlineKeyboardButtonData("📅 Повторы", "orders_repeat"),
			tgbotapi.NewInlineKeyboardButtonData("💸 На кассу", "orders_cash"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = keyboard
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Failed to send orders menu: %v", err)
	}
}
func (h *TelegramHandler) handleCallback(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	// 1. Удаляем сообщение целиком
	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := h.bot.Request(del); err != nil {
		log.Printf("Failed to delete message: %v", err)
	}
	// 2. Обработка команд
	switch {
	case query.Data == "orders_today":
		h.showTodayOrders(chatID)

	case query.Data == "orders_repeat":
		h.showRepeatOrders(chatID)

	case query.Data == "orders_cash":
		h.showCashOrders(chatID)

	case query.Data == "init":
		h.showMainMenu(chatID)

	case strings.HasPrefix(query.Data, "accept_happy_"):
		h.handleAcceptOrder(query, "happy")

	case strings.HasPrefix(query.Data, "accept_"):
		h.handleAcceptOrder(query, "normal")

	case strings.HasPrefix(query.Data, "order_view_"):
		h.showOrderDetails(query)
	}

	// 3. Подтверждаем обработку callback
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := h.bot.Request(callback); err != nil {
		log.Printf("Callback error: %v", err)
	}
}
