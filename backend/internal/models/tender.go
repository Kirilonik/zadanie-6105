package models

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;table:tender"`
	Name            string    `gorm:"type:varchar(100);not null"`
	Description     string    `gorm:"type:text"`
	ServiceType     string    `gorm:"type:varchar(50)"`
	Status          string    `gorm:"type:varchar(20)"`
	OrganizationID  uuid.UUID `gorm:"type:uuid;not null"`
	CreatorUsername string    `gorm:"type:varchar(50);not null"`
	Version         int       `gorm:"type:integer;not null;default:1"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type TenderHistory struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;table:tender_history"`
	TenderID        uuid.UUID `gorm:"type:uuid;not null"`
	Name            string    `gorm:"type:varchar(100);not null"`
	Description     string    `gorm:"type:text"`
	ServiceType     string    `gorm:"type:varchar(50)"`
	Status          string    `gorm:"type:varchar(20)"`
	OrganizationID  uuid.UUID `gorm:"type:uuid;not null"`
	CreatorUsername string    `gorm:"type:varchar(50);not null"`
	Version         int       `gorm:"type:integer;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
