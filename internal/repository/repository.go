package repository

import (
	"context"
	"http_server/internal/domain"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) error
	GetById(ctx context.Context, id uuid.UUID) (domain.Task, error)
	List(ctx context.Context) ([]domain.Task, error)
	ListByCompletedStatus(ctx context.Context, status bool) ([]domain.Task, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetVersionById(ctx context.Context, id uuid.UUID) (int, error)
}
