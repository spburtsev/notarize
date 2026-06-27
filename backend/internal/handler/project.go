package handler

import (
	"context"

	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/oas"
)

func (h *ServerHandler) CreateProject(ctx context.Context, req *oas.CreateProjectRequest) (*oas.Project, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}
	project := models.Project{Name: req.Name, CreatedBy: principal.UserID}
	if err := h.db.WithContext(ctx).Create(&project).Error; err != nil {
		return nil, err
	}
	out := toOASProject(project)
	return &out, nil
}

func (h *ServerHandler) ListProjects(ctx context.Context, params oas.ListProjectsParams) (*oas.ProjectPage, error) {
	q := func() *gorm.DB { return h.db.WithContext(ctx).Model(&models.Project{}) }
	var total int64
	if err := q().Count(&total).Error; err != nil {
		return nil, err
	}
	var projects []models.Project
	if err := q().Order("created_at DESC").Offset(listOffset(params.Offset)).Limit(listLimit(params.Limit)).Find(&projects).Error; err != nil {
		return nil, err
	}
	items := make([]oas.Project, len(projects))
	for i, p := range projects {
		items[i] = toOASProject(p)
	}
	return &oas.ProjectPage{Items: items, Total: int32(total)}, nil
}

func (h *ServerHandler) GetProject(ctx context.Context, params oas.GetProjectParams) (*oas.Project, error) {
	var project models.Project
	if err := h.db.WithContext(ctx).First(&project, "id = ?", params.ProjectId).Error; err != nil {
		return nil, notFound(err)
	}
	out := toOASProject(project)
	return &out, nil
}

func toOASProject(p models.Project) oas.Project {
	return oas.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedBy: p.CreatedBy,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
