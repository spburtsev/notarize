package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/spburtsev/notarize/internal/db/models"
)

func TestHashCompare(t *testing.T) {
	h, err := Hash("s3cret")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !Compare(h, "s3cret") {
		t.Fatal("expected matching password to compare true")
	}
	if Compare(h, "wrong") {
		t.Fatal("expected wrong password to compare false")
	}
}

func TestIssueParse(t *testing.T) {
	svc := NewService("test-secret", time.Hour)
	u := models.User{ID: uuid.New(), Role: models.UserRoleManager}

	tok, err := svc.Issue(u)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	p, err := svc.parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if p.UserID != u.ID {
		t.Fatalf("user id: got %v want %v", p.UserID, u.ID)
	}
	if p.Role != u.Role {
		t.Fatalf("role: got %v want %v", p.Role, u.Role)
	}
}

func TestParseRejectsExpiredAndTampered(t *testing.T) {
	u := models.User{ID: uuid.New(), Role: models.UserRoleAdmin}

	expired := NewService("test-secret", -time.Hour)
	tok, _ := expired.Issue(u)
	if _, err := expired.parse(tok); err == nil {
		t.Fatal("expected expired token to be rejected")
	}

	signed, _ := NewService("secret-a", time.Hour).Issue(u)
	if _, err := NewService("secret-b", time.Hour).parse(signed); err == nil {
		t.Fatal("expected token signed with a different secret to be rejected")
	}
}
