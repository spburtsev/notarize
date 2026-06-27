package handler

import (
	"context"

	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/oas"
)

func (h *ServerHandler) CreateUser(ctx context.Context, req *oas.CreateUserRequest) (*oas.User, error) {
	hash, err := auth.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: hash,
		Role:         models.UserRoleEmployee,
	}
	if v, ok := req.Role.Get(); ok {
		user.Role = models.UserRole(v)
	}
	if err := h.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	out := toOASUser(user)
	return &out, nil
}

func (h *ServerHandler) ListUsers(ctx context.Context, params oas.ListUsersParams) (*oas.UserPage, error) {
	q := func() *gorm.DB {
		base := h.db.WithContext(ctx).Model(&models.User{})
		if v, ok := params.Role.Get(); ok {
			base = base.Where("role = ?", string(v))
		}
		return base
	}
	var total int64
	if err := q().Count(&total).Error; err != nil {
		return nil, err
	}
	var users []models.User
	if err := q().Order("created_at DESC").Offset(listOffset(params.Offset)).Limit(listLimit(params.Limit)).Find(&users).Error; err != nil {
		return nil, err
	}
	items := make([]oas.User, len(users))
	for i, u := range users {
		items[i] = toOASUser(u)
	}
	return &oas.UserPage{Items: items, Total: int32(total)}, nil
}
