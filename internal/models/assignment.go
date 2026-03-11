package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Title string

	Duration int

	CreatedAt time.Time
}

func (a *Assignment) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}
