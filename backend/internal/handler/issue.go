package handler

import (
	"context"

	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/oas"
)

func (h *ServerHandler) CreateIssue(ctx context.Context, req *oas.CreateIssueRequest) (*oas.Issue, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}
	issue := models.Issue{
		ProjectID: req.ProjectID,
		Title:     req.Title,
		Status:    models.IssueStatusOpen,
		CreatedBy: principal.UserID,
	}
	if v, ok := req.Description.Get(); ok {
		issue.Description = &v
	}
	if err := h.db.WithContext(ctx).Create(&issue).Error; err != nil {
		return nil, err
	}
	out := toOASIssue(issue)
	return &out, nil
}

func (h *ServerHandler) ListIssues(ctx context.Context, params oas.ListIssuesParams) (*oas.IssuePage, error) {
	q := func() *gorm.DB {
		base := h.db.WithContext(ctx).Model(&models.Issue{})
		if v, ok := params.ProjectId.Get(); ok {
			base = base.Where("project_id = ?", v)
		}
		if v, ok := params.Status.Get(); ok {
			base = base.Where("status = ?", string(v))
		}
		return base
	}
	var total int64
	if err := q().Count(&total).Error; err != nil {
		return nil, err
	}
	var issues []models.Issue
	if err := q().Order("created_at DESC").Offset(listOffset(params.Offset)).Limit(listLimit(params.Limit)).Find(&issues).Error; err != nil {
		return nil, err
	}
	items := make([]oas.Issue, len(issues))
	for i, is := range issues {
		items[i] = toOASIssue(is)
	}
	return &oas.IssuePage{Items: items, Total: int32(total)}, nil
}

func (h *ServerHandler) GetIssue(ctx context.Context, params oas.GetIssueParams) (*oas.Issue, error) {
	var issue models.Issue
	if err := h.db.WithContext(ctx).First(&issue, "id = ?", params.IssueId).Error; err != nil {
		return nil, notFound(err)
	}
	out := toOASIssue(issue)
	return &out, nil
}

func toOASIssue(i models.Issue) oas.Issue {
	out := oas.Issue{
		ID:        i.ID,
		ProjectID: i.ProjectID,
		Title:     i.Title,
		Status:    oas.IssueStatus(i.Status),
		CreatedBy: i.CreatedBy,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
	if i.Description != nil {
		out.Description = oas.NewOptNilString(*i.Description)
	}
	return out
}
