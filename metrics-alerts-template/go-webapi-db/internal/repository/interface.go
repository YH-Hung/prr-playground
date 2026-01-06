package repository

import (
	"context"

	"go-webapi-db/internal/model"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error
	FindByStatus(ctx context.Context, status string) ([]*model.User, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

