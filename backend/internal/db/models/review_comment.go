package models

import (
	"time"

	"github.com/google/uuid"
)

type ReviewComment struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProcessID uuid.UUID    `gorm:"type:uuid;not null;index"`
	StepID    *uuid.UUID   `gorm:"type:uuid"`
	Step      *ProcessStep `gorm:"foreignKey:StepID"`
	AuthorID  uuid.UUID    `gorm:"type:uuid;not null"`
	Body      string       `gorm:"type:text;not null"`
	CreatedAt time.Time
}
