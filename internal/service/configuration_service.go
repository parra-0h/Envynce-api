package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type ConfigurationService struct {
	repo    *repository.ConfigurationRepository
	appRepo *repository.ApplicationRepository
	envRepo *repository.EnvironmentRepository
}

func NewConfigurationService(repo *repository.ConfigurationRepository, appRepo *repository.ApplicationRepository, envRepo *repository.EnvironmentRepository) *ConfigurationService {
	return &ConfigurationService{
		repo:    repo,
		appRepo: appRepo,
		envRepo: envRepo,
	}
}

func (s *ConfigurationService) CreateConfiguration(ctx context.Context, config *domain.Configuration) error {
	// Verify app and env exist
	_, err := s.appRepo.GetByID(ctx, config.ApplicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}
	_, err = s.envRepo.GetByID(ctx, config.EnvironmentID)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	// Check if a previous version exists
	latest, err := s.repo.GetLatest(ctx, config.ApplicationID, config.EnvironmentID, config.Key)
	if err == nil && latest != nil {
		// Archive previous version if it was active
		if latest.Status == "active" {
			err = s.repo.UpdateStatus(ctx, latest.ID, "archived")
			if err != nil {
				return fmt.Errorf("failed to archive previous version: %w", err)
			}
		}
		config.Version = latest.Version + 1
	} else {
		config.Version = 1
	}

	config.Status = "active"
	err = s.repo.Create(ctx, config)
	if err != nil {
		return err
	}

	// Audit Log
	newVal, _ := json.Marshal(config)
	var oldVal string
	if latest != nil {
		ov, _ := json.Marshal(latest)
		oldVal = string(ov)
	}

	s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		Action:   "CREATE/UPDATE",
		Entity:   "Configuration",
		EntityID: config.ID,
		OldValue: oldVal,
		NewValue: string(newVal),
	})

	return nil
}

func (s *ConfigurationService) GetActiveConfigs(ctx context.Context, appID, envID uint) ([]domain.Configuration, error) {
	return s.repo.GetAll(ctx, appID, envID)
}

func (s *ConfigurationService) GetByID(ctx context.Context, id uint) (*domain.Configuration, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ConfigurationService) GetAuditLogs(ctx context.Context) ([]domain.AuditLog, error) {
	return s.repo.GetAuditLogs(ctx)
}
