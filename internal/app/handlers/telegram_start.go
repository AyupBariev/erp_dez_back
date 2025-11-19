package handlers

import (
	"erp/internal/app/models"
	"erp/internal/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func (h *TelegramHandler) handleStart(msg *tgbotapi.Message) {
	telegramID := msg.From.ID
	username := msg.From.UserName

	// 1️⃣ Пытаемся найти инженера по Telegram ID
	engineer, err := h.engineerService.GetEngineerByTelegramID(telegramID)
	if err != nil {
		h.sendMessage(msg.Chat.ID, "Ошибка при проверке пользователя. Попробуйте позже.")
		return
	}

	// 2️⃣ Если не найден по Telegram ID — пробуем по username
	if engineer == nil && username != "" {
		engineer, err = h.engineerService.GetEngineerByUsername(username)
		if err != nil {
			h.sendMessage(msg.Chat.ID, "Ошибка при проверке пользователя. Попробуйте позже.")
			return
		}

		// Если нашли по username, но без Telegram ID — обновляем
		if engineer != nil && engineer.GetTelegramID() == nil {
			err = h.engineerService.UpdateTelegramID(int64(engineer.ID), telegramID)
			if err != nil {
				h.sendMessage(msg.Chat.ID, "Ошибка при привязке Telegram ID. Попробуйте позже.")
				return
			}
			h.sendMessage(msg.Chat.ID, "Ваш Telegram успешно привязан к профилю инженера.")
		}
	}

	// 3️⃣ Если не найден вообще — создаём нового
	if engineer == nil {
		newEngineer := &models.Engineer{
			Username:   username,
			TelegramID: utils.Int64ToNullInt64(telegramID),
			IsApproved: false,
		}

		engineer, err = h.engineerService.CreateEngineer(newEngineer)
		if err != nil {
			h.sendMessage(msg.Chat.ID, "Ошибка при регистрации. Попробуйте позже.")
			return
		}

		h.sendMessage(msg.Chat.ID, "Вы зарегистрированы. Ожидайте подтверждения администратора.")
		return
	}

	// 4️⃣ Проверяем статус подтверждения
	if !engineer.IsApproved {
		h.sendMessage(msg.Chat.ID, "Ваш профиль ожидает подтверждения администратора.")
		return
	}

	// 5️⃣ Если всё ок — показываем меню
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
	if err := h.orderService.EngineerAcceptOrderByErpNumber(erpNumber); err != nil {
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

	note := order.Note
	if note == "" {
		note = "—"
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
			"📅 Дата и время: %s\n"+
			"📞 Телефон: %s\n"+
			"💰Сумма: %s\n",
		order.ERPNumber,
		clientName,
		address,
		note,
		date,
		phoneDisplay,
		price,
	)

	reportLink, _ := h.reportService.GenerateReportLink(int64(order.ID), int64(order.Engineer.ID))

	// 🎛 Кнопки под заказом
	var buttons []tgbotapi.InlineKeyboardButton

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("🏠 В меню", "init"))

	if reportLink != "" {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonURL("📄 Отправить отчёт", reportLink))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	message, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке деталей заказа: %v", err)
	}
	go func(chatID int64, msgID int) {
		time.Sleep(1 * time.Minute)
		del := tgbotapi.NewDeleteMessage(chatID, msgID)
		if _, err := h.bot.Request(del); err != nil {
			log.Print("del err:", err)
		}
	}(query.Message.Chat.ID, message.MessageID)
}

func parseErpNumber(data string) int64 {
	parts := strings.Split(data, "_")
	if len(parts) < 2 {
		return 0
	}
	num, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	return num
}
