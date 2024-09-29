package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"size:20;not null;unique" json:"username"`
	Password string    `gorm:"size:60;not null" json:"password"`
	RoleID   uint      `gorm:"not null" json:"-"`
	Role     Role      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"role"`
}

func (User) TableName() string {
	return "users"
}
