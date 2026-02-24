package repository

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type ConfigurationRepository struct {
	db *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Create(ctx context.Context, config *domain.Configuration) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *ConfigurationRepository) GetLatest(ctx context.Context, appID, envID uint, key string) (*domain.Configuration, error) {
	var config domain.Configuration
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND environment_id = ? AND key = ?", appID, envID, key).
		Order("version DESC").
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) GetAll(ctx context.Context, appID, envID uint) ([]domain.Configuration, error) {
	var configs []domain.Configuration
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND environment_id = ? AND status = 'active'", appID, envID).
		Find(&configs).Error
	return configs, err
}

func (r *ConfigurationRepository) GetByID(ctx context.Context, id uint) (*domain.Configuration, error) {
	var config domain.Configuration
	err := r.db.WithContext(ctx).Preload("Application").Preload("Environment").First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Configuration{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ConfigurationRepository) GetAuditLogs(ctx context.Context) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *ConfigurationRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
