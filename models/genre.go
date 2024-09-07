package models

type Genre struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:20;not null" json:"name"`
}

func (Genre) TableName() string {
	return "genres"
}
