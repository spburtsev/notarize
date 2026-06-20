package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleManager  UserRole = "MANAGER"
	UserRoleEmployee UserRole = "EMPLOYEE"
	UserRoleInvitee  UserRole = "INVITEE"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"type:text;not null;unique"`
	FirstName    string    `gorm:"type:text;not null"`
	LastName     string    `gorm:"type:text;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	Role         UserRole  `gorm:"type:varchar(32);not null;default:EMPLOYEE;check:users_role_check,role IN ('ADMIN','MANAGER','EMPLOYEE','INVITEE')"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
