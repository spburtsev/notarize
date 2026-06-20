package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:text;not null"`
	Issues    []Issue   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
