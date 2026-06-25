package models

import (
	"time"

	"github.com/google/uuid"
)

type ProcessingStatus string

const (
	ProcessingStatusQueued     ProcessingStatus = "QUEUED"
	ProcessingStatusProcessing ProcessingStatus = "PROCESSING"
	ProcessingStatusDone       ProcessingStatus = "DONE"
	ProcessingStatusFailed     ProcessingStatus = "FAILED"
)

type DocumentProcessingJob struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DocumentID uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex"`
	Document   *Document        `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE"`
	Status     ProcessingStatus `gorm:"type:varchar(16);not null;default:QUEUED;index;check:document_processing_jobs_status_check,status IN ('QUEUED','PROCESSING','DONE','FAILED')"`
	Attempts   int32            `gorm:"not null;default:0"`
	Error      *string          `gorm:"type:text"`
	StartedAt  *time.Time
	FinishedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
