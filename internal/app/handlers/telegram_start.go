package handlers

import (
	"erp/internal/app/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func (h *TelegramHandler) handleStart(msg *tgbotapi.Message) {
	// 1. Пробуем найти инженера
	engineer, err := h.engineerService.GetEngineerByTelegramID(msg.From.ID)
	if err != nil {
		h.sendMessage(msg.Chat.ID, "Ошибка при проверке пользователя. Попробуйте позже.")
		return
	}

	// 2. Если нет — создаём нового с минимальными данными
	if engineer == nil {
		engineer = &models.Engineer{
			Username:   msg.From.UserName,
			TelegramID: msg.From.ID,
		}
		engineer, err = h.engineerService.CreateEngineer(engineer)
		if err != nil {
			h.sendMessage(msg.Chat.ID, "Ошибка при регистрации. Попробуйте позже.")
			return
		}
	}

	// 3. Если инженер не подтверждён
	if !engineer.IsApproved {
		h.sendMessage(msg.Chat.ID, "Ваш запрос отправлен администратору. Ожидайте подтверждения.")
		return
	}

	// 4. Если подтверждён — показываем меню
	h.showMainMenu(msg.Chat.ID)
}

func (h *TelegramHandler) showMainMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Главное меню")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Мои заказы"),
		),
	)
	_, err := h.bot.Send(msg)
	if err != nil {
		return
	}
}

func (h *TelegramHandler) handleAcceptOrder(query *tgbotapi.CallbackQuery, acceptMode string) {
	erpNumber := parseErpNumber(query.Data)
	if erpNumber == 0 {
		h.sendMessage(query.Message.Chat.ID, "Ошибка: не удалось определить заказ")
		return
	}

	// 🧹 Удаляем сообщение логиста
	deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	if _, err := h.bot.Request(deleteMsg); err != nil {
		log.Printf("❌ Не удалось удалить сообщение: %v", err)
	}

	// ✅ Обновляем статус заказа
	if err := h.orderService.MarkAsAcceptedByErpNumber(erpNumber); err != nil {
		log.Printf("Ошибка обновления заказа: %v", err)
		h.sendMessage(query.Message.Chat.ID, "Ошибка при обновлении заказа 😕")
		return
	}

	// 💬 Формируем сообщение
	var text string
	if acceptMode == "happy" {
		text = fmt.Sprintf("😄 Вы с удовольствием приняли заказ №%d!", erpNumber)
	} else {
		text = fmt.Sprintf("✅ Вы приняли заказ №%d.", erpNumber)
	}

	// 🧭 Добавляем inline-кнопки для дальнейших действий
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Посмотреть заказ", fmt.Sprintf("order_view_%d", erpNumber)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Вернуться в меню", "init"),
		),
	)

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("❌ Не удалось отправить сообщение инженеру: %v", err)
	}
}

func (h *TelegramHandler) showOrderDetails(query *tgbotapi.CallbackQuery) {
	// Извлекаем ERP номер из callback data
	parts := strings.Split(query.Data, "_")
	if len(parts) < 3 {
		h.sendMessage(query.Message.Chat.ID, "Ошибка: не удалось определить заказ 😕")
		return
	}

	erpNumber, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		h.sendMessage(query.Message.Chat.ID, "Ошибка: неверный номер заказа 😕")
		return
	}

	// Получаем заказ из БД
	order, err := h.orderService.GetOrderForAssign(erpNumber)
	if err != nil || order == nil {
		h.sendMessage(query.Message.Chat.ID, fmt.Sprintf("Заказ №%d не найден 🕵️‍♂️", erpNumber))
		return
	}

	// 🧾 Формируем сообщение
	clientName := order.ClientName
	if clientName == "" {
		clientName = "—"
	}

	address := order.Address
	if address == "" {
		address = "—"
	}

	problem := order.Problem
	if problem == "" {
		problem = "—"
	}

	date := ""
	if !order.ScheduledAt.IsZero() {
		date = order.ScheduledAt.Format("02.01.2006 15:04")
	} else {
		date = "не указано"
	}

	text := fmt.Sprintf(
		"📄 *Информация о заказе №%d*\n\n"+
			"👤 Клиент: %s\n"+
			"🏠 Адрес: %s\n"+
			"🔧 Проблема: %s\n"+
			"📅 Дата и время: %s\n",
		order.ERPNumber,
		clientName,
		address,
		problem,
		date,
	)

	// 🎛 Кнопки под заказом
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 В меню", "init"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Повторить заказ", fmt.Sprintf("order_repeat_%d", order.ERPNumber)),
		),
	)

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке деталей заказа: %v", err)
	}
}

func parseErpNumber(data string) int64 {
	parts := strings.Split(data, "_")
	if len(parts) < 2 {
		return 0
	}
	num, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	return num
}
