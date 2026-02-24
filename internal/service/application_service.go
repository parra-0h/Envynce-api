package service

import (
	"context"

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
	return s.repo.GetByID(ctx, id)
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
