package repository

import (
	"context"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *ApplicationRepository) GetAll(ctx context.Context) ([]domain.Application, error) {
	var apps []domain.Application
	err := r.db.WithContext(ctx).Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) GetByName(ctx context.Context, name string) (*domain.Application, error) {
	var app domain.Application
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id uint) (*domain.Application, error) {
	var app domain.Application
	err := r.db.WithContext(ctx).First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) Update(ctx context.Context, app *domain.Application) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *ApplicationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Application{}, id).Error
}

func (r *ApplicationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Application{}).Count(&count).Error
	return count, err
}
