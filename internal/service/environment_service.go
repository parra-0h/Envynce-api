package service

import (
	"context"
	"errors"

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
	env, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("environment not found")
	}
	return env, nil
}

func (s *EnvironmentService) UpdateEnvironment(ctx context.Context, id uint, name, description string) (*domain.Environment, error) {
	env, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("environment not found")
	}
	if name != "" {
		env.Name = name
	}
	if description != "" {
		env.Description = description
	}
	if err := s.repo.Update(ctx, env); err != nil {
		return nil, errors.New("failed to update environment")
	}
	return env, nil
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("environment not found")
	}
	return s.repo.Delete(ctx, id)
}
