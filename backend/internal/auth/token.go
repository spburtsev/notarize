package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/spburtsev/notarize/internal/db/models"
)

type claims struct {
	Role models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) Issue(user models.User) (string, error) {
	now := time.Now()
	c := claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.secret)
}

func (s *Service) parse(token string) (Principal, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorized
		}
		return s.secret, nil
	})
	if err != nil {
		return Principal{}, ErrUnauthorized
	}
	id, err := uuid.Parse(c.Subject)
	if err != nil {
		return Principal{}, ErrUnauthorized
	}
	return Principal{UserID: id, Role: c.Role}, nil
}
