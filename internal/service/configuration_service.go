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
	auditRepo   *repository.AuditLogRepository
}

func NewConfigurationService(
	repo *repository.ConfigurationRepository,
	appRepo *repository.ApplicationRepository,
	envRepo *repository.EnvironmentRepository,
	versionRepo *repository.ConfigVersionRepository,
	auditRepo *repository.AuditLogRepository,
) *ConfigurationService {
	return &ConfigurationService{
		repo:        repo,
		appRepo:     appRepo,
		envRepo:     envRepo,
		versionRepo: versionRepo,
		auditRepo:   auditRepo,
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
	config.CreatedBy = userName
	if err := s.repo.Create(ctx, config); err != nil {
		return errors.New("failed to create configuration")
	}

	// Save version snapshot
	_ = s.versionRepo.Create(ctx, &domain.ConfigVersion{
		ConfigurationID: config.ID,
		Key:             config.Key,
		Value:           config.Value,
		Description:     config.Description,
		Active:          true,
		VersionNumber:   config.Version,
	})

	// Audit Log
	_ = s.auditRepo.LogAction(ctx, userID, "CREATE", "Configuration", config.ID, config)

	return nil
}

func (s *ConfigurationService) GetConfigsAsMap(ctx context.Context, appName, envName string) (map[string]string, uint, uint, error) {
	app, err := s.appRepo.GetByName(ctx, appName)
	if err != nil {
		return nil, 0, 0, errors.New("application not found")
	}
	env, err := s.envRepo.GetByName(ctx, envName)
	if err != nil {
		return nil, 0, 0, errors.New("environment not found")
	}

	configs, err := s.repo.GetAll(ctx, app.ID, env.ID)
	if err != nil {
		return nil, 0, 0, err
	}

	result := make(map[string]string)
	for _, c := range configs {
		result[c.Key] = c.Value
	}
	return result, app.ID, env.ID, nil
}

func (s *ConfigurationService) GetActiveConfigs(ctx context.Context, appID, envID uint) ([]domain.Configuration, error) {
	return s.repo.GetAll(ctx, appID, envID)
}

func (s *ConfigurationService) GetConfigurationHistory(ctx context.Context, id uint) ([]domain.Configuration, error) {
	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetHistory(ctx, config.ApplicationID, config.EnvironmentID, config.Key)
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
		CreatedBy:     userName,
	}
	if err := s.repo.Create(ctx, newConfig); err != nil {
		return nil, errors.New("failed to create new version")
	}

	// Save version snapshot
	_ = s.versionRepo.Create(ctx, &domain.ConfigVersion{
		ConfigurationID: newConfig.ID,
		Key:             newConfig.Key,
		Value:           newConfig.Value,
		Description:     newConfig.Description,
		Active:          true,
		VersionNumber:   newConfig.Version,
	})

	// Audit Log
	_ = s.auditRepo.LogAction(ctx, userID, "UPDATE", "Configuration", newConfig.ID, newConfig)

	return newConfig, nil
}

func (s *ConfigurationService) DeleteConfiguration(ctx context.Context, id uint) error {
	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("configuration not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Audit Log
	_ = s.auditRepo.LogAction(ctx, 0, "DELETE", "Configuration", id, config)
	return nil
}

func (s *ConfigurationService) SearchConfigurations(ctx context.Context, appID, envID uint, keyword string) ([]domain.Configuration, error) {
	return s.repo.Search(ctx, appID, envID, keyword)
}

func (s *ConfigurationService) GetConfigVersions(ctx context.Context, configID uint) ([]domain.ConfigVersion, error) {
	return s.versionRepo.GetByConfigID(ctx, configID)
}

func (s *ConfigurationService) GetAuditLogs(ctx context.Context) ([]domain.AuditLog, error) {
	return s.auditRepo.GetAll(ctx, 50)
}

func (s *ConfigurationService) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	totalApps, _ := s.appRepo.Count(ctx)
	totalEnvs, _ := s.envRepo.Count(ctx)
	totalConfigs, _ := s.repo.Count(ctx)
	activeConfigs, _ := s.repo.CountByStatus(ctx, "active")
	recentUpdates, _ := s.auditRepo.GetAll(ctx, 5)

	return &domain.DashboardStats{
		TotalApps:         totalApps,
		TotalConfigs:      totalConfigs,
		ActiveConfigs:     activeConfigs,
		TotalEnvironments: totalEnvs,
		RecentUpdates:     recentUpdates,
	}, nil
}
func (s *ConfigurationService) RestoreVersion(ctx context.Context, versionID uint, userID uint, userName string) (*domain.Configuration, error) {
	// 1. Get the version to restore
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, errors.New("version history entry not found")
	}

	// 2. Get the current active configuration for this key/app/env
	// We need to find the configuration that this version belongs to.
	// Actually, version.ConfigurationID points to A specific configuration record.
	// But we want the LATEST record for that Key/App/Env to archive it.
	config, err := s.repo.GetByID(ctx, version.ConfigurationID)
	if err != nil {
		return nil, errors.New("parent configuration not found")
	}

	latest, err := s.repo.GetLatest(ctx, config.ApplicationID, config.EnvironmentID, config.Key)
	if err != nil {
		return nil, errors.New("failed to find current active version")
	}

	// 3. Archive current latest
	if latest.Status == "active" {
		if err := s.repo.UpdateStatus(ctx, latest.ID, "archived"); err != nil {
			return nil, errors.New("failed to archive current version")
		}
	}

	// 4. Create new version based on history
	newConfig := &domain.Configuration{
		Key:           version.Key,
		Value:         version.Value,
		Description:   version.Description,
		ApplicationID: config.ApplicationID,
		EnvironmentID: config.EnvironmentID,
		Version:       latest.Version + 1,
		Status:        "active",
		CreatedBy:     userName,
	}
	if err := s.repo.Create(ctx, newConfig); err != nil {
		return nil, errors.New("failed to create restored version")
	}

	// 5. Save version snapshot
	_ = s.versionRepo.Create(ctx, &domain.ConfigVersion{
		ConfigurationID: newConfig.ID,
		Key:             newConfig.Key,
		Value:           newConfig.Value,
		Description:     newConfig.Description,
		Active:          true,
		VersionNumber:   newConfig.Version,
	})

	// 6. Audit Log
	_ = s.auditRepo.LogAction(ctx, userID, "RESTORE", "Configuration", newConfig.ID, newConfig)

	return newConfig, nil
}
