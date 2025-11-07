package order

type Repository interface {
	GetOrders(date *string) ([]*Order, error)
	Create(order *Order) error
	Update(order *Order) error
	Delete(erpNumber int64) error
	GetNextERPNumber() (int64, error)
	GetByERPNumber(erpNumber int64) (*Order, error)
	GetTodayOrders(chatID int64) ([]Order, error)
	GetRepeatOrders(chatID int64) ([]Order, error)
	GetCashOrders(chatID int64) ([]Order, error)
}
