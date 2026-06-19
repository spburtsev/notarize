package handler

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spburtsev/notarize/internal/db/sqlc"
	"github.com/spburtsev/notarize/internal/oas"
	"github.com/spburtsev/notarize/internal/storage"
)

type ServerHandler struct {
	pool    *pgxpool.Pool
	q       *sqlc.Queries
	storage *storage.Storage
}

func New(pool *pgxpool.Pool, store *storage.Storage) *ServerHandler {
	return &ServerHandler{
		pool:    pool,
		q:       sqlc.New(pool),
		storage: store,
	}
}

func (h *ServerHandler) CancelApprovalProcess(ctx context.Context, params oas.CancelApprovalProcessParams) (*oas.ApprovalProcess, error) {
	return nil, nil
}
func (h *ServerHandler) CreateApprovalProcess(ctx context.Context, req *oas.CreateApprovalProcessRequest) (*oas.ApprovalProcess, error) {
	return nil, nil
}
func (h *ServerHandler) CreateDocument(ctx context.Context, req *oas.CreateDocumentReq) (*oas.Document, error) {
	return nil, nil
}
func (h *ServerHandler) CreateFolder(ctx context.Context, req *oas.CreateFolderRequest) (*oas.Folder, error) {
	return nil, nil
}
func (h *ServerHandler) CreateProcessComment(ctx context.Context, req *oas.CreateReviewCommentRequest, params oas.CreateProcessCommentParams) (*oas.ReviewComment, error) {
	return nil, nil
}
func (h *ServerHandler) DeleteFolder(ctx context.Context, params oas.DeleteFolderParams) error {
	return nil
}
func (h *ServerHandler) DownloadDocumentContent(ctx context.Context, params oas.DownloadDocumentContentParams) (oas.DownloadDocumentContentOK, error) {
	return oas.DownloadDocumentContentOK{}, nil
}
func (h *ServerHandler) GetApprovalProcess(ctx context.Context, params oas.GetApprovalProcessParams) (*oas.ApprovalProcess, error) {
	return nil, nil
}
func (h *ServerHandler) GetDocument(ctx context.Context, params oas.GetDocumentParams) (*oas.Document, error) {
	return nil, nil
}
func (h *ServerHandler) GetFolder(ctx context.Context, params oas.GetFolderParams) (*oas.Folder, error) {
	return nil, nil
}
func (h *ServerHandler) ListApprovalProcesses(ctx context.Context, params oas.ListApprovalProcessesParams) (*oas.ApprovalProcessPage, error) {
	return nil, nil
}
func (h *ServerHandler) ListDocuments(ctx context.Context, params oas.ListDocumentsParams) (*oas.DocumentPage, error) {
	return nil, nil
}
func (h *ServerHandler) ListFolders(ctx context.Context, params oas.ListFoldersParams) (*oas.FolderPage, error) {
	return nil, nil
}
func (h *ServerHandler) ListProcessComments(ctx context.Context, params oas.ListProcessCommentsParams) (*oas.ReviewCommentPage, error) {
	return nil, nil
}
func (h *ServerHandler) MoveDocument(ctx context.Context, req *oas.MoveDocumentRequest, params oas.MoveDocumentParams) (*oas.Document, error) {
	return nil, nil
}
func (h *ServerHandler) SubmitDecision(ctx context.Context, req *oas.SubmitDecisionRequest, params oas.SubmitDecisionParams) (*oas.Decision, error) {
	return nil, nil
}
func (h *ServerHandler) UpdateFolder(ctx context.Context, req *oas.UpdateFolderRequest, params oas.UpdateFolderParams) (*oas.Folder, error) {
	return nil, nil
}

func (h *ServerHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	return &oas.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: oas.Error{
			ErrorCode: "internal_server_error",
			Message:   err.Error(),
		},
	}
}
