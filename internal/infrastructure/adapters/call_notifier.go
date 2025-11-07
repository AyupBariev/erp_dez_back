package adapters

import (
	"erp/internal/app/services"
	"erp/internal/domain/order"
)

type CallNotifier struct {
	service *services.CallService
}

func NewCallNotifier(service *services.CallService) *CallNotifier {
	return &CallNotifier{service: service}
}

func (c *CallNotifier) NotifyEngineerNewOrder(o *order.Order) {
	//if o.Engineer == nil {
	//	return
	//}
	//
	//// здесь логика звонка инженеру
	//c.service.ScheduleEngineerCall(o.Engineer.ID, o)
}
