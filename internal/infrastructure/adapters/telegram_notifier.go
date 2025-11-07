package adapters

import (
	"erp/internal/app/services"
	"erp/internal/domain/order"
)

type TelegramNotifier struct {
	service *services.TelegramService
}

func NewTelegramNotifier(service *services.TelegramService) *TelegramNotifier {
	return &TelegramNotifier{service: service}
}

func (t *TelegramNotifier) NotifyEngineerNewOrder(o *order.Order) {
	if o.Engineer == nil || o.Engineer.TelegramID == nil {
		return
	}

	text := order.FormatOrderMessage(o)
	t.service.SendMessageWithKeyboard(*o.Engineer.TelegramID, text, order.OrderInlineKeyboard(o))
}
