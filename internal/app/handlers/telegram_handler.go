package handlers

import (
	"erp/internal/app/models"
	"erp/internal/app/services"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramHandler struct {
	bot             *tgbotapi.BotAPI
	engineerService *services.EngineerService
	orderService    *services.OrderService
	reportService   *services.ReportService
	keyboards       map[string]tgbotapi.InlineKeyboardMarkup
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, engineerService *services.EngineerService, orderService *services.OrderService, reportService *services.ReportService) *TelegramHandler {
	return &TelegramHandler{
		bot:             bot,
		engineerService: engineerService,
		orderService:    orderService,
		reportService:   reportService,
		keyboards: map[string]tgbotapi.InlineKeyboardMarkup{
			"init": tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🤖 Запустить бота", "init"),
				),
			),
		},
	}
}

func (h *TelegramHandler) HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				h.handleStart(update.Message)
			}
			continue
		}

		if update.Message != nil {
			switch update.Message.Text {
			case "Мои заказы":
				h.showOrdersMenu(update.Message.Chat.ID)
			case "Сегодня":
				h.showTodayOrders(update.Message.Chat.ID)
			}
			continue
		}

		if update.CallbackQuery != nil {
			h.handleCallback(update.CallbackQuery)
		}
	}
}

func (h *TelegramHandler) sendMessage(chatID int64, text string, keyboardKey ...string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if len(keyboardKey) > 0 {
		if kb, ok := h.keyboards[keyboardKey[0]]; ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (h *TelegramHandler) sendOrderList(chatID int64, title string, orders []models.Order) {
	text := title + "\n"
	for _, order := range orders {
		text += fmt.Sprintf("- %d (%s)\n", order.ERPNumber, order.Status)
	}
	h.sendMessage(chatID, text)
}

func (h *TelegramHandler) showRepeatOrders(chatID int64) {
	orders, err := h.orderService.GetRepeatOrders(chatID) // Теперь получаем 2 значения
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения повторяющихся заказов")
		log.Printf("GetRepeatOrders error: %v", err)
		return
	}

	if len(orders) == 0 {
		h.sendMessage(chatID, "Нет повторяющихся заказов")
		return
	}

	h.sendOrderList(chatID, "Повторяющиеся заказы:", orders)
}

func (h *TelegramHandler) showCashOrders(chatID int64) {
	orders, err := h.orderService.GetCashOrders(chatID) // Теперь получаем 2 значения
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения заказов на кассу")
		log.Printf("GetCashOrders error: %v", err)
		return
	}

	if len(orders) == 0 {
		h.sendMessage(chatID, "Нет заказов на кассу")
		return
	}

	h.sendOrderList(chatID, "Заказы на кассу:", orders)
}

func (h *TelegramHandler) showTodayOrders(chatID int64) {
	orders, err := h.orderService.GetTodayOrders(chatID)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения сегодняшних заказов")
		log.Printf("GetTodayOrders error: %v", err)
		return
	}

	if len(orders) == 0 {
		h.sendMessage(chatID, "На сегодня заказов нет")
		return
	}

	// Заголовок
	msg := tgbotapi.NewMessage(chatID, "📅 *Заказы на сегодня:*")
	msg.ParseMode = "Markdown"

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, order := range orders {
		// Формируем время (если указано)
		timeStr := "Без времени"
		if !order.ScheduledAt.IsZero() {
			timeStr = order.ScheduledAt.Format("15:04")
		}

		// Формируем текст кнопки
		btnText := fmt.Sprintf("%s. %s. %s (Открыть заказ)", timeStr, order.Problem.Name, order.Address)

		// Создаём кнопку с callback data вида "order_view_100002"
		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("order_view_%d", order.ERPNumber))

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	// Добавим кнопку "Назад в меню"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 Вернуться в меню", "init"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("send orders list error: %v", err)
	}
}
