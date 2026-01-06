package repository

import (
	"context"
	"time"

	"go-webapi-db/internal/model"
	"go-webapi-db/internal/metrics"
)

const (
	appName  = "go-webapi-db"
	database = "go_webapi_db"
)

// InstrumentedUserRepository wraps UserRepository with metrics collection
type InstrumentedUserRepository struct {
	repo *UserRepository
}

// NewInstrumentedUserRepository creates a new instrumented repository
func NewInstrumentedUserRepository(repo *UserRepository) *InstrumentedUserRepository {
	return &InstrumentedUserRepository{
		repo: repo,
	}
}

// Ensure InstrumentedUserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*InstrumentedUserRepository)(nil)

func (r *InstrumentedUserRepository) Create(ctx context.Context, user *model.User) error {
	start := time.Now()
	err := r.repo.Create(ctx, user)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "insert", "users", duration, err)
	return err
}

func (r *InstrumentedUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	start := time.Now()
	user, err := r.repo.FindByID(ctx, id)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "find", "users", duration, err)
	return user, err
}

func (r *InstrumentedUserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	start := time.Now()
	users, err := r.repo.FindAll(ctx)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "find", "users", duration, err)
	return users, err
}

func (r *InstrumentedUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	start := time.Now()
	user, err := r.repo.FindByEmail(ctx, email)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "find", "users", duration, err)
	return user, err
}

func (r *InstrumentedUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	start := time.Now()
	exists, err := r.repo.ExistsByEmail(ctx, email)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "count", "users", duration, err)
	return exists, err
}

func (r *InstrumentedUserRepository) Update(ctx context.Context, id string, user *model.User) error {
	start := time.Now()
	err := r.repo.Update(ctx, id, user)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "update", "users", duration, err)
	return err
}

func (r *InstrumentedUserRepository) Delete(ctx context.Context, id string) error {
	start := time.Now()
	err := r.repo.Delete(ctx, id)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "delete", "users", duration, err)
	return err
}

func (r *InstrumentedUserRepository) FindByStatus(ctx context.Context, status string) ([]*model.User, error) {
	start := time.Now()
	users, err := r.repo.FindByStatus(ctx, status)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "find", "users", duration, err)
	return users, err
}

func (r *InstrumentedUserRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	start := time.Now()
	count, err := r.repo.CountByStatus(ctx, status)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "count", "users", duration, err)
	return count, err
}

func (r *InstrumentedUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	start := time.Now()
	exists, err := r.repo.ExistsByID(ctx, id)
	duration := time.Since(start)
	
	metrics.RecordOperation(appName, database, "count", "users", duration, err)
	return exists, err
}

