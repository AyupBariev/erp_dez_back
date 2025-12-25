package models

type Problem struct {
	ID   int64  `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name"`
}

func (Problem) TableName() string {
	return "problems"
}
