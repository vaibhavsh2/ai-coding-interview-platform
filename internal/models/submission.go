package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submission struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	QuestionID    uuid.UUID `gorm:"type:uuid;not null"`
	Language      string
	SourceCode    string
	Status        string
	PassedTests   int
	TotalTests    int
	ExecutionTime int
	MemoryUsed    int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Submission) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
