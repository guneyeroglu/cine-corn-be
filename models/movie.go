package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name           string    `gorm:"size:50;not null" json:"name"`
	Point          string    `gorm:"size:4;not null" json:"point"`
	ThumbnailImage string    `gorm:"not null" json:"thumbnailImage"`
}

type MovieDetails struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Point       string    `gorm:"size:4;not null" json:"point"`
	Description string    `gorm:"type:text;not null" json:"description"`
	ReleaseDate time.Time `gorm:"type:date;not null" json:"releaseDate"`
	RunTime     string    `gorm:"size:3;not null" json:"runTime"`
	GenreIDs    []uint    `gorm:"-" json:"-"`
	Genres      []string  `gorm:"-" json:"genres"`
	PosterImage string    `gorm:"not null" json:"posterImage"`
	HeroImage   string    `gorm:"not null" json:"heroImage"`
	IsFeatured  bool      `gorm:"not null" json:"isFeatured"`
	IsNew       bool      `gorm:"not null" json:"isNew"`
}
