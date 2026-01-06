package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"go-webapi-db/internal/model"
	"go-webapi-db/internal/repository"
)

type UserService struct {
	repo          repository.UserRepositoryInterface
	metrics       *MetricsService
	random        *rand.Rand
}

func NewUserService(repo repository.UserRepositoryInterface, metrics *MetricsService) *UserService {
	return &UserService{
		repo:    repo,
		metrics: metrics,
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("create_failed")
		return nil, err
	}
	if exists {
		s.metrics.IncrementUserOperationErrors("duplicate_email")
		return nil, errors.New("user with email " + req.Email + " already exists")
	}

	// Simulate occasional slow operations
	s.simulateRandomDelay()

	user := &model.User{
		Email:  req.Email,
		Name:   req.Name,
		Status: "ACTIVE",
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("create_failed")
		return nil, err
	}

	s.metrics.IncrementUserCreated()
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	// Simulate occasional timeout scenarios
	if s.random.Intn(100) < 2 { // 2% chance
		time.Sleep(3 * time.Second)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("timeout")
		return nil, err
	}

	if user == nil {
		s.metrics.IncrementUserOperationErrors("not_found")
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	s.simulateRandomDelay()
	return s.repo.FindAll(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		s.metrics.IncrementUserOperationErrors("not_found")
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("update_failed")
		return nil, err
	}

	if user == nil {
		s.metrics.IncrementUserOperationErrors("not_found")
		return nil, errors.New("user not found with id: " + id)
	}

	// Simulate occasional errors
	if s.random.Intn(100) < 1 { // 1% chance
		s.metrics.IncrementUserOperationErrors("update_failed")
		return nil, errors.New("simulated database error")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	err = s.repo.Update(ctx, id, user)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("update_failed")
		return nil, err
	}

	// Fetch updated user
	updatedUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.metrics.IncrementUserUpdated()
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("delete_failed")
		return err
	}

	if !exists {
		s.metrics.IncrementUserOperationErrors("not_found")
		return errors.New("user not found")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperationErrors("delete_failed")
		return err
	}

	s.metrics.IncrementUserDeleted()
	return nil
}

func (s *UserService) GetUsersByStatus(ctx context.Context, status string) ([]*model.User, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	return s.repo.FindByStatus(ctx, status)
}

func (s *UserService) CountUsersByStatus(ctx context.Context, status string) (int64, error) {
	stopTimer := s.metrics.StartUserOperationTimer()
	defer stopTimer()

	return s.repo.CountByStatus(ctx, status)
}

func (s *UserService) CallExternalService(ctx context.Context, serviceName string) (string, error) {
	start := time.Now()
	
	// Simulate external call
	delay := time.Duration(s.random.Intn(500)+100) * time.Millisecond
	time.Sleep(delay)

	// Simulate occasional failures
	if s.random.Intn(100) < 5 { // 5% failure rate
		duration := time.Since(start)
		s.metrics.RecordExternalCallDuration(serviceName, duration)
		s.metrics.IncrementExternalCallErrors(serviceName)
		return "", errors.New("external service " + serviceName + " failed")
	}

	duration := time.Since(start)
	s.metrics.RecordExternalCallDuration(serviceName, duration)
	return "Success from " + serviceName, nil
}

func (s *UserService) simulateRandomDelay() {
	delay := time.Duration(s.random.Intn(190)+10) * time.Millisecond
	time.Sleep(delay)
}

