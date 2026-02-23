package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CodingQuestion struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string
	Description string
	Difficulty  string
	TimeLimit   int
	MemoryLimit int

	TestCases []TestCase `gorm:"foreignKey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *CodingQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	return
}
