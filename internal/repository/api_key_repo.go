package repository

import (
	"context"
	"time"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

func (r *APIKeyRepository) GetAll(ctx context.Context) ([]domain.APIKey, error) {
	var keys []domain.APIKey
	err := r.db.WithContext(ctx).Preload("Applications").Order("created_at DESC").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) FindByHashedKey(ctx context.Context, hashedKey string) (*domain.APIKey, error) {
	var apiKey domain.APIKey
	err := r.db.WithContext(ctx).Preload("Applications").Where("hashed_key = ? AND status = 'active'", hashedKey).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.APIKey{}).Where("id = ?", id).Update("last_used", &now).Error
}

func (r *APIKeyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.APIKey{}, id).Error
}

func (r *APIKeyRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.APIKey{}).Where("id = ?", id).Update("status", status).Error
}

func (r *APIKeyRepository) GetByID(ctx context.Context, id uint) (*domain.APIKey, error) {
	var apiKey domain.APIKey
	err := r.db.WithContext(ctx).Preload("Applications").First(&apiKey, id).Error
	return &apiKey, err
}
