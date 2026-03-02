package repository

import (
	"context"
	"strings"

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
	query := r.db.WithContext(ctx).Where("status = 'active'")
	if appID > 0 {
		query = query.Where("application_id = ?", appID)
	}
	if envID > 0 {
		query = query.Where("environment_id = ?", envID)
	}
	err := query.Preload("Application").Preload("Environment").Find(&configs).Error
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

func (r *ConfigurationRepository) Update(ctx context.Context, config *domain.Configuration) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *ConfigurationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Configuration{}, id).Error
}

func (r *ConfigurationRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Configuration{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ConfigurationRepository) Search(ctx context.Context, appID, envID uint, keyword string) ([]domain.Configuration, error) {
	var configs []domain.Configuration
	keyword = strings.ToLower(keyword)
	query := r.db.WithContext(ctx).Where("status = 'active'")
	if appID > 0 {
		query = query.Where("application_id = ?", appID)
	}
	if envID > 0 {
		query = query.Where("environment_id = ?", envID)
	}
	if keyword != "" {
		query = query.Where("LOWER(key) LIKE ? OR LOWER(value) LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := query.Preload("Application").Preload("Environment").Find(&configs).Error
	return configs, err
}
