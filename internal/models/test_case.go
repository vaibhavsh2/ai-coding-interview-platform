package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestCase struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	QuestionID     uuid.UUID `gorm:"type:uuid;not null"`
	Input          string
	ExpectedOutput string
	IsHidden       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (t *TestCase) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
