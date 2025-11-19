package domain

type AggregatorPayoutRepository interface {
	GetDayPayouts(from, to string) ([]*AggregatorDayPayout, error)
	MarkPaid(channel string, date string, amount float64) error
}
