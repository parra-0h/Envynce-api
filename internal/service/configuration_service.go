package service

import (
	"context"
	"errors"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type ConfigurationService struct {
	repo        *repository.ConfigurationRepository
	appRepo     *repository.ApplicationRepository
	envRepo     *repository.EnvironmentRepository
	versionRepo *repository.ConfigVersionRepository
}

func NewConfigurationService(
	repo *repository.ConfigurationRepository,
	appRepo *repository.ApplicationRepository,
	envRepo *repository.EnvironmentRepository,
	versionRepo *repository.ConfigVersionRepository,
) *ConfigurationService {
	return &ConfigurationService{
		repo:        repo,
		appRepo:     appRepo,
		envRepo:     envRepo,
		versionRepo: versionRepo,
	}
}

func (s *ConfigurationService) CreateConfiguration(ctx context.Context, config *domain.Configuration, userID uint, userName string) error {
	// Verify app and env exist
	_, err := s.appRepo.GetByID(ctx, config.ApplicationID)
	if err != nil {
		return errors.New("application not found")
	}
	_, err = s.envRepo.GetByID(ctx, config.EnvironmentID)
	if err != nil {
		return errors.New("environment not found")
	}

	// Check if a previous version exists and archive it
	latest, err := s.repo.GetLatest(ctx, config.ApplicationID, config.EnvironmentID, config.Key)
	if err == nil && latest != nil {
		if latest.Status == "active" {
			if err := s.repo.UpdateStatus(ctx, latest.ID, "archived"); err != nil {
				return errors.New("failed to archive previous version")
			}
		}
		config.Version = latest.Version + 1
	} else {
		config.Version = 1
	}

	config.Status = "active"
	if err := s.repo.Create(ctx, config); err != nil {
		return errors.New("failed to create configuration")
	}

	// Save version snapshot
	s.versionRepo.Create(ctx, &domain.ConfigVersion{
		ConfigurationID: config.ID,
		Key:             config.Key,
		Value:           config.Value,
		Version:         config.Version,
		ChangedByUserID: userID,
		ChangedByName:   userName,
	})

	return nil
}

func (s *ConfigurationService) GetActiveConfigs(ctx context.Context, appID, envID uint) ([]domain.Configuration, error) {
	return s.repo.GetAll(ctx, appID, envID)
}

func (s *ConfigurationService) GetByID(ctx context.Context, id uint) (*domain.Configuration, error) {
	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("configuration not found")
	}
	return config, nil
}

func (s *ConfigurationService) UpdateConfiguration(ctx context.Context, id uint, req *domain.UpdateConfigRequest, userID uint, userName string) (*domain.Configuration, error) {
	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("configuration not found")
	}

	// Archive current version
	if err := s.repo.UpdateStatus(ctx, id, "archived"); err != nil {
		return nil, errors.New("failed to archive current version")
	}

	// Create new version
	newConfig := &domain.Configuration{
		Key:           config.Key,
		Value:         req.Value,
		Description:   req.Description,
		ApplicationID: config.ApplicationID,
		EnvironmentID: config.EnvironmentID,
		Version:       config.Version + 1,
		Status:        "active",
	}
	if err := s.repo.Create(ctx, newConfig); err != nil {
		return nil, errors.New("failed to create new version")
	}

	// Save version snapshot
	s.versionRepo.Create(ctx, &domain.ConfigVersion{
		ConfigurationID: newConfig.ID,
		Key:             newConfig.Key,
		Value:           newConfig.Value,
		Version:         newConfig.Version,
		ChangedByUserID: userID,
		ChangedByName:   userName,
	})

	return newConfig, nil
}

func (s *ConfigurationService) DeleteConfiguration(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("configuration not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *ConfigurationService) SearchConfigurations(ctx context.Context, appID, envID uint, keyword string) ([]domain.Configuration, error) {
	return s.repo.Search(ctx, appID, envID, keyword)
}

func (s *ConfigurationService) GetConfigVersions(ctx context.Context, configID uint) ([]domain.ConfigVersion, error) {
	return s.versionRepo.GetByConfigID(ctx, configID)
}
