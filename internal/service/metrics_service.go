package service

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type MetricsService struct {
	repo *repository.RequestLogRepository
}

func NewMetricsService(repo *repository.RequestLogRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) GetRequestsPerMinute(ctx context.Context) ([]domain.MetricResponse, error) {
	return s.repo.GetRequestsPerMinute(ctx, true)
}

func (s *MetricsService) LogRequest(ctx context.Context, apiKeyID, appID, envID uint) error {
	log := &domain.RequestLog{
		APIKeyID:      apiKeyID,
		ApplicationID: appID,
		EnvironmentID: envID,
	}
	return s.repo.Create(ctx, log)
}
