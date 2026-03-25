package models

import (
	"time"

	"github.com/google/uuid"
)

type CandidateAssignment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	CandidateID  uuid.UUID `gorm:"type:uuid;not null;index"`
	AssignmentID string    `gorm:"type:varchar;not null;index"` // Matching UUID/VARCHAR of Assignment table
	Status       string    `gorm:"default:'ASSIGNED'"` // ASSIGNED, IN_PROGRESS, COMPLETED
	CreatedAt    time.Time
}
