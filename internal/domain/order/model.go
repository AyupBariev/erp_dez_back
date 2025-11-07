package order

import "time"

type Status string

type Order struct {
	ID                uint
	ERPNumber         int64
	AggregatorID      int64
	OurPercent        float64
	Price             string
	FinishPrice       string
	ClientName        string
	EngineerID        *int64
	AdminID           int64
	Phones            []string
	Address           string
	WorkVolume        string
	ProblemID         *int64
	Note              string
	ScheduledAt       time.Time
	Status            string
	ConfirmedAt       *time.Time
	RepeatID          int
	RepeatDescription string
	RepeatedBy        string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
