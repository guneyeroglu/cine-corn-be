package models

import "github.com/google/uuid"

type UserFavoriteMovie struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	MovieID uuid.UUID `gorm:"type:uuid;not null" json:"movieId"`
}

func (UserFavoriteMovie) TableName() string {
	return "user_favorite_movies"
}
