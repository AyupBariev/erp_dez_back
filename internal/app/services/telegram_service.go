package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

type TelegramService struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramService(bot *tgbotapi.BotAPI) *TelegramService {
	return &TelegramService{Bot: bot}
}

func (s *TelegramService) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := s.Bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения Telegram: %v", err)
	}
}

func (s *TelegramService) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := s.Bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (s *TelegramService) sendMessageWithTTL(chatID int64, text string, ttl time.Duration) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// кнопка «Удалить»
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", "del_"+time.Now().Add(ttl).Format(time.RFC3339)),
		),
	)

	sent, err := s.Bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
		return
	}

	// удалить сообщение через 1 минуту
	go func(chatID int64, msgID int) {
		time.Sleep(time.Minute)
		del := tgbotapi.NewDeleteMessage(chatID, msgID)
		if _, err := s.Bot.Request(del); err != nil {
			log.Print("del err:", err)
		}
	}(chatID, sent.MessageID)

}
