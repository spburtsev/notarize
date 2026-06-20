package models

import (
	"time"

	"github.com/google/uuid"
)

type DocumentStatus string

const (
	DocumentStatusUploaded DocumentStatus = "UPLOADED"
	DocumentStatusInReview DocumentStatus = "IN_REVIEW"
	DocumentStatusApproved DocumentStatus = "APPROVED"
	DocumentStatusRejected DocumentStatus = "REJECTED"
	DocumentStatusArchived DocumentStatus = "ARCHIVED"
)

type Document struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string            `gorm:"type:text;not null"`
	ContentType string            `gorm:"type:text;not null"`
	SizeBytes   int64             `gorm:"not null"`
	SHA256      string            `gorm:"column:sha256;type:text;not null"`
	IssueID     uuid.UUID         `gorm:"type:uuid;not null;index"`
	Issue       *Issue            `gorm:"foreignKey:IssueID"`
	Status      DocumentStatus    `gorm:"type:varchar(32);not null;default:UPLOADED;index;check:documents_status_check,status IN ('UPLOADED','IN_REVIEW','APPROVED','REJECTED','ARCHIVED')"`
	StorageKey  string            `gorm:"type:text;not null"`
	Processes   []ApprovalProcess `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE"`
	CreatedBy   uuid.UUID         `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
