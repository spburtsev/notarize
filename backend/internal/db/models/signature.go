package models

import (
	"time"

	"github.com/google/uuid"
)

type SignatureAlgorithm string

const (
	SignatureAlgorithmED25519         SignatureAlgorithm = "ED25519"
	SignatureAlgorithmECDSAP256SHA256 SignatureAlgorithm = "ECDSA_P256_SHA256"
	SignatureAlgorithmRSAPSSSHA256    SignatureAlgorithm = "RSA_PSS_SHA256"
)

type Signature struct {
	ID                  uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DecisionID          uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex"`
	Algorithm           SignatureAlgorithm `gorm:"type:varchar(32);not null;check:signatures_algorithm_check,algorithm IN ('ED25519','ECDSA_P256_SHA256','RSA_PSS_SHA256')"`
	Value               string             `gorm:"type:text;not null"`
	KeyID               string             `gorm:"column:key_id;type:text;not null"`
	Certificate         *string            `gorm:"type:text"`
	SignedPayloadSHA256 string             `gorm:"column:signed_payload_sha256;type:text;not null"`
	SignedAt            time.Time          `gorm:"not null"`
}
