package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null" json:"name"`
}

func (Role) TableName() string {
	return "roles"
}
