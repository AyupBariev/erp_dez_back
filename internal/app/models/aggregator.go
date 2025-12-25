package models

type Aggregator struct {
	ID   int64  `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name"`
}

func (Aggregator) TableName() string {
	return "aggregators"
}
