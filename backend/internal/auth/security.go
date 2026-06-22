package auth

import (
	"context"
	"slices"

	"github.com/spburtsev/notarize/internal/oas"
)

func (s *Service) HandleBearerAuth(ctx context.Context, _ oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	p, err := s.parse(t.Token)
	if err != nil {
		return ctx, err
	}
	if len(t.Roles) > 0 && !slices.Contains(t.Roles, string(p.Role)) {
		return ctx, ErrForbidden
	}
	return withPrincipal(ctx, p), nil
}
