package service

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type EnvironmentService struct {
	repo *repository.EnvironmentRepository
}

func NewEnvironmentService(repo *repository.EnvironmentRepository) *EnvironmentService {
	return &EnvironmentService{repo: repo}
}

func (s *EnvironmentService) CreateEnvironment(ctx context.Context, env *domain.Environment) error {
	return s.repo.Create(ctx, env)
}

func (s *EnvironmentService) GetAllEnvironments(ctx context.Context) ([]domain.Environment, error) {
	return s.repo.GetAll(ctx)
}

func (s *EnvironmentService) GetEnvironmentByID(ctx context.Context, id uint) (*domain.Environment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
