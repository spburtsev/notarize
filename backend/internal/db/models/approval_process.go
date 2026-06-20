package models

import (
	"time"

	"github.com/google/uuid"
)

type ProcessStatus string

const (
	ProcessStatusPending    ProcessStatus = "PENDING"
	ProcessStatusInProgress ProcessStatus = "IN_PROGRESS"
	ProcessStatusApproved   ProcessStatus = "APPROVED"
	ProcessStatusRejected   ProcessStatus = "REJECTED"
	ProcessStatusCancelled  ProcessStatus = "CANCELLED"
)

type ApprovalProcess struct {
	ID               uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DocumentID       uuid.UUID     `gorm:"type:uuid;not null;index"`
	Document         *Document     `gorm:"foreignKey:DocumentID"`
	Status           ProcessStatus `gorm:"type:varchar(32);not null;default:PENDING;index;check:approval_processes_status_check,status IN ('PENDING','IN_PROGRESS','APPROVED','REJECTED','CANCELLED')"`
	CurrentStepIndex *int32
	Steps            []ProcessStep   `gorm:"foreignKey:ProcessID;constraint:OnDelete:CASCADE"`
	Comments         []ReviewComment `gorm:"foreignKey:ProcessID;constraint:OnDelete:CASCADE"`
	CreatedBy        uuid.UUID       `gorm:"type:uuid;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
