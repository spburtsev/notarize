package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type StepStatus string

const (
	StepStatusPending  StepStatus = "PENDING"
	StepStatusActive   StepStatus = "ACTIVE"
	StepStatusApproved StepStatus = "APPROVED"
	StepStatusRejected StepStatus = "REJECTED"
	StepStatusSkipped  StepStatus = "SKIPPED"
)

type StepPolicy string

const (
	StepPolicyAny    StepPolicy = "ANY"
	StepPolicyAll    StepPolicy = "ALL"
	StepPolicyQuorum StepPolicy = "QUORUM"
)

type ProcessStep struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProcessID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_process_steps_process_index,priority:1"`
	Index             int32      `gorm:"not null;uniqueIndex:idx_process_steps_process_index,priority:2"`
	Name              string     `gorm:"type:text;not null"`
	Policy            StepPolicy `gorm:"type:varchar(16);not null;check:process_steps_policy_check,policy IN ('ANY','ALL','QUORUM')"`
	RequiredApprovals *int32
	Status            StepStatus     `gorm:"type:varchar(16);not null;default:PENDING;check:process_steps_status_check,status IN ('PENDING','ACTIVE','APPROVED','REJECTED','SKIPPED')"`
	ApproverUserIDs   pq.StringArray `gorm:"column:approver_user_ids;type:uuid[];not null;default:'{}'"`
	ApproverRoles     pq.StringArray `gorm:"column:approver_roles;type:text[];not null;default:'{}'"`
	Decisions         []Decision     `gorm:"foreignKey:StepID;constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
