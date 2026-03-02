package service

import (
	"context"
	"errors"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, app *domain.Application) error {
	return s.repo.Create(ctx, app)
}

func (s *ApplicationService) GetAllApplications(ctx context.Context) ([]domain.Application, error) {
	return s.repo.GetAll(ctx)
}

func (s *ApplicationService) GetApplicationByID(ctx context.Context, id uint) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("application not found")
	}
	return app, nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, id uint, name, description string) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("application not found")
	}
	if name != "" {
		app.Name = name
	}
	if description != "" {
		app.Description = description
	}
	if err := s.repo.Update(ctx, app); err != nil {
		return nil, errors.New("failed to update application")
	}
	return app, nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("application not found")
	}
	return s.repo.Delete(ctx, id)
}
