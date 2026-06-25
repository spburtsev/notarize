package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"

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

	contentType := req.File.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
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
