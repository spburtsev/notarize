package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/oas"
)

func (h *ServerHandler) CreateDocument(ctx context.Context, req *oas.CreateDocumentReq) (*oas.Document, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}

	b, err := io.ReadAll(req.File.File)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(b)
	if contentType != "application/pdf" {
		return nil, ErrNotPDF
	}

	name := req.File.Name
	if req.Name.Set && req.Name.Value != "" {
		name = req.Name.Value
	}

	digest := sha256.Sum256(b)
	id := uuid.New()
	storageKey := "documents/" + id.String()

	if err := h.storage.Put(ctx, storageKey, bytes.NewReader(b), int64(len(b)), contentType); err != nil {
		return nil, err
	}

	doc := models.Document{
		ID:          id,
		Name:        name,
		ContentType: contentType,
		SizeBytes:   int64(len(b)),
		SHA256:      hex.EncodeToString(digest[:]),
		IssueID:     req.IssueID,
		Status:      models.DocumentStatusUploaded,
		StorageKey:  storageKey,
		CreatedBy:   principal.UserID,
	}

	if err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}
		return tx.Create(&models.DocumentProcessingJob{DocumentID: doc.ID}).Error
	}); err != nil {
		return nil, err
	}

	out := toOASDocument(doc)
	return &out, nil
}

func (h *ServerHandler) ListDocuments(ctx context.Context, params oas.ListDocumentsParams) (*oas.DocumentPage, error) {
	q := h.db.WithContext(ctx).Order("created_at DESC").Limit(listLimit(params.Limit))
	if v, ok := params.IssueId.Get(); ok {
		q = q.Where("issue_id = ?", v)
	}
	if v, ok := params.Status.Get(); ok {
		q = q.Where("status = ?", string(v))
	}
	var docs []models.Document
	if err := q.Find(&docs).Error; err != nil {
		return nil, err
	}
	items := make([]oas.Document, len(docs))
	for i, d := range docs {
		items[i] = toOASDocument(d)
	}
	return &oas.DocumentPage{Items: items}, nil
}

func (h *ServerHandler) DownloadDocumentContent(ctx context.Context, params oas.DownloadDocumentContentParams) (oas.DownloadDocumentContentOK, error) {
	var doc models.Document
	if err := h.db.WithContext(ctx).First(&doc, "id = ?", params.DocumentId).Error; err != nil {
		return oas.DownloadDocumentContentOK{}, notFound(err)
	}

	obj, err := h.storage.Get(ctx, doc.StorageKey)
	if err != nil {
		return oas.DownloadDocumentContentOK{}, err
	}
	return oas.DownloadDocumentContentOK{Data: obj}, nil
}

func (h *ServerHandler) GetDocumentResult(ctx context.Context, params oas.GetDocumentResultParams) (*oas.DocumentResult, error) {
	var job models.DocumentProcessingJob
	if err := h.db.WithContext(ctx).First(&job, "document_id = ?", params.DocumentId).Error; err != nil {
		return nil, notFound(err)
	}

	out := oas.DocumentResult{Status: oas.ProcessingStatus(job.Status)}
	if job.Error != nil {
		out.Error = oas.OptNilString{Value: *job.Error, Set: true}
	}
	if job.Status == models.ProcessingStatusDone {
		obj, err := h.storage.Get(ctx, "derived/"+params.DocumentId.String()+".md")
		if err != nil {
			return nil, err
		}
		defer obj.Close()
		b, err := io.ReadAll(obj)
		if err != nil {
			return nil, err
		}
		out.Output = oas.OptNilString{Value: string(b), Set: true}
	}
	return &out, nil
}

func (h *ServerHandler) ReprocessDocument(ctx context.Context, params oas.ReprocessDocumentParams) (*oas.DocumentResult, error) {
	res := h.db.WithContext(ctx).Model(&models.DocumentProcessingJob{}).
		Where("document_id = ?", params.DocumentId).
		Updates(map[string]any{
			"status":      models.ProcessingStatusQueued,
			"error":       nil,
			"started_at":  nil,
			"finished_at": nil,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return &oas.DocumentResult{Status: oas.ProcessingStatusQUEUED}, nil
}

func toOASDocument(d models.Document) oas.Document {
	return oas.Document{
		ID:          d.ID,
		Name:        d.Name,
		ContentType: d.ContentType,
		SizeBytes:   d.SizeBytes,
		SHA256:      d.SHA256,
		IssueID:     d.IssueID,
		Status:      oas.DocumentStatus(d.Status),
		CreatedBy:   d.CreatedBy,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
