package repository

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) Create(ctx context.Context, env *domain.Environment) error {
	return r.db.WithContext(ctx).Create(env).Error
}

func (r *EnvironmentRepository) GetAll(ctx context.Context) ([]domain.Environment, error) {
	var envs []domain.Environment
	err := r.db.WithContext(ctx).Find(&envs).Error
	return envs, err
}

func (r *EnvironmentRepository) GetByID(ctx context.Context, id uint) (*domain.Environment, error) {
	var env domain.Environment
	err := r.db.WithContext(ctx).First(&env, id).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvironmentRepository) Update(ctx context.Context, env *domain.Environment) error {
	return r.db.WithContext(ctx).Save(env).Error
}

func (r *EnvironmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Environment{}, id).Error
}

func (r *EnvironmentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Environment{}).Count(&count).Error
	return count, err
}
