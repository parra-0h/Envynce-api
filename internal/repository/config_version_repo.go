package repository

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type ConfigVersionRepository struct {
	db *gorm.DB
}

func NewConfigVersionRepository(db *gorm.DB) *ConfigVersionRepository {
	return &ConfigVersionRepository{db: db}
}

func (r *ConfigVersionRepository) Create(ctx context.Context, v *domain.ConfigVersion) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *ConfigVersionRepository) GetByConfigID(ctx context.Context, configID uint) ([]domain.ConfigVersion, error) {
	var versions []domain.ConfigVersion
	err := r.db.WithContext(ctx).
		Where("configuration_id = ?", configID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}
