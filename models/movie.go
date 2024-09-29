package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Movie struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name          string    `gorm:"size:50;not null" json:"name"`
	Point         string    `gorm:"size:4;not null" json:"point"`
	PosterImage   string    `gorm:"type:text;not null" json:"posterImage"`
	IsFavorite    bool      `gorm:"-" json:"isFavorite"`
	IsAddedToList bool      `gorm:"-" json:"isAddedToList"`
}

func (Movie) TableName() string {
	return "movies"
}

type MovieDetails struct {
	Movie
	BannerImage string         `gorm:"not null" json:"bannerImage"`
	Description string         `gorm:"type:text;not null" json:"description"`
	ReleaseDate time.Time      `gorm:"type:date;not null" json:"releaseDate"`
	RunTime     string         `gorm:"size:3;not null" json:"runTime"`
	Genres      []Genre        `gorm:"many2many:movie_genres;foreignKey:ID;joinForeignKey:MovieID;References:ID;joinReferences:GenreID" json:"-"`
	GenreNames  []string       `gorm:"-" json:"genres"`
	IsFeatured  bool           `gorm:"not null" json:"isFeatured"`
	IsNew       bool           `gorm:"not null" json:"isNew"`
	Stars       pq.StringArray `gorm:"type:text[];not null;" json:"stars"`
}

func (MovieDetails) TableName() string {
	return "movies"
}

type MovieDetailsRequest struct {
	ID uuid.UUID `json:"id"`
}
