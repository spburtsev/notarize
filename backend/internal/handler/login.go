package handler

import (
	"context"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/oas"
)

func (h *ServerHandler) Login(ctx context.Context, req *oas.LoginRequest) (*oas.LoginResponse, error) {
	var user models.User
	err := h.db.WithContext(ctx).Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, auth.ErrUnauthorized
	}
	if !auth.Compare(user.PasswordHash, req.Password) {
		return nil, auth.ErrUnauthorized
	}

	token, err := h.auth.Issue(user)
	if err != nil {
		return nil, err
	}
	return &oas.LoginResponse{Token: token, User: toOASUser(user)}, nil
}

func toOASUser(u models.User) oas.User {
	return oas.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      oas.UserRole(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
