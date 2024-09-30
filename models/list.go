package models

import "github.com/google/uuid"

type UserListMovie struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	MovieID uuid.UUID `gorm:"type:uuid;not null" json:"movieId"`
}

func (UserListMovie) TableName() string {
	return "user_list_movies"
}

type UserListMovieRequest struct {
	MovieID uuid.UUID `gorm:"type:uuid;not null" json:"movieId"`
}

func (UserListMovieRequest) TableName() string {
	return "user_list_movies"
}
