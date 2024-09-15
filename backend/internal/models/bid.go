package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	TenderID    uuid.UUID `gorm:"type:uuid;not null"`
	AuthorType  string    `gorm:"type:varchar(20);not null"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null"`
	Version     int       `gorm:"default:1"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type BidReview struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	BidID       uuid.UUID `gorm:"type:uuid;not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type BidHistory struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	BidID       uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	TenderID    uuid.UUID `gorm:"type:uuid;not null"`
	AuthorType  string    `gorm:"type:varchar(20);not null"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null"`
	Version     int       `gorm:"default:1"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
