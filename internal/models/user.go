package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"unique"`
	Password  string
	Name      string
	Provider  string // local, google
	CreatedAt time.Time
}
