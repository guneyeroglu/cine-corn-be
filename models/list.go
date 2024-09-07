package models

import "github.com/google/uuid"

type List struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"size:50;not null" json:"name"`
	MovieIds []uint    `gorm:"-" json:"-"`
	Movies   []Movie   `gorm:"-" json:"movies"`
}

func (List) TableName() string {
	return "lists"
}
