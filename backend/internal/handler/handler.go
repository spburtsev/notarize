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

var ErrNotFound = errors.New("not found")

var ErrNotPDF = errors.New("not a pdf")

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func listLimit(opt oas.OptInt32) int {
	if v, ok := opt.Get(); ok && v > 0 {
		if v > 200 {
			return 200
		}
		return int(v)
	}
	return 50
}

func listOffset(opt oas.OptInt32) int {
	if v, ok := opt.Get(); ok && v > 0 {
		return int(v)
	}
	return 0
}

func (h *ServerHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	switch {
	case errors.Is(err, ErrNotPDF):
		return errResponse(http.StatusBadRequest, "invalid_document", "only PDF documents are accepted")
	case errors.Is(err, ErrNotFound):
		return errResponse(http.StatusNotFound, "not_found", "the requested resource was not found")
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
