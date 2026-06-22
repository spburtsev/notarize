package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/oas"
	"github.com/spburtsev/notarize/internal/storage"
)

type ServerHandler struct {
	oas.UnimplementedHandler

	db      *gorm.DB
	storage *storage.Storage
	auth    *auth.Service
}

func New(db *gorm.DB, store *storage.Storage, authSvc *auth.Service) *ServerHandler {
	return &ServerHandler{db: db, storage: store, auth: authSvc}
}

func (h *ServerHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	switch {
	case errors.Is(err, auth.ErrForbidden):
		return errResponse(http.StatusForbidden, "forbidden", "you do not have permission to perform this action")
	case errors.Is(err, auth.ErrUnauthorized):
		return errResponse(http.StatusUnauthorized, "unauthorized", "authentication required")
	}
	if _, ok := errors.AsType[*ogenerrors.SecurityError](err); ok {
		return errResponse(http.StatusUnauthorized, "unauthorized", "authentication required")
	}
	return errResponse(http.StatusInternalServerError, "internal_server_error", err.Error())
}

func errResponse(status int, code, message string) *oas.ErrorStatusCode {
	return &oas.ErrorStatusCode{
		StatusCode: status,
		Response:   oas.Error{ErrorCode: code, Message: message},
	}
}
