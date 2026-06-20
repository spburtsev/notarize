package models

import (
	"time"

	"github.com/google/uuid"
)

type IssueStatus string

const (
	IssueStatusOpen       IssueStatus = "OPEN"
	IssueStatusInProgress IssueStatus = "IN_PROGRESS"
	IssueStatusClosed     IssueStatus = "CLOSED"
)

type Issue struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID   uuid.UUID   `gorm:"type:uuid;not null;index"`
	Project     *Project    `gorm:"foreignKey:ProjectID"`
	Title       string      `gorm:"type:text;not null"`
	Description *string     `gorm:"type:text"`
	Status      IssueStatus `gorm:"type:varchar(16);not null;default:OPEN;index;check:issues_status_check,status IN ('OPEN','IN_PROGRESS','CLOSED')"`
	Documents   []Document  `gorm:"foreignKey:IssueID;constraint:OnDelete:CASCADE"`
	CreatedBy   uuid.UUID   `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
