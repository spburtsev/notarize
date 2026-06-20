package models

import (
	"time"

	"github.com/google/uuid"
)

type DecisionType string

const (
	DecisionTypeApprove DecisionType = "APPROVE"
	DecisionTypeReject  DecisionType = "REJECT"
)

type Decision struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StepID    uuid.UUID    `gorm:"type:uuid;not null;index"`
	ActorID   uuid.UUID    `gorm:"type:uuid;not null"`
	Type      DecisionType `gorm:"type:varchar(16);not null;check:decisions_type_check,type IN ('APPROVE','REJECT')"`
	Comment   *string      `gorm:"type:text"`
	Signature *Signature   `gorm:"foreignKey:DecisionID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
}
