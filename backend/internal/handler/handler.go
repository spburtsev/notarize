package handler

import (
	"context"
	"net/http"

	"github.com/spburtsev/notarize/internal/oas"
)

// ServerHandler implements the generated oas.Handler interface. Endpoint logic
// is not wired up yet: embedding UnimplementedHandler supplies a "not
// implemented" response for every operation until each one is filled in.
type ServerHandler struct {
	oas.UnimplementedHandler
}

// NewError maps an otherwise-unhandled error to a 500 response.
func (h *ServerHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	return &oas.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: oas.Error{
			ErrorCode: "internal_server_error",
			Message:   err.Error(),
		},
	}
}
