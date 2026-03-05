package repository

import (
	"context"
	"time"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type RequestLogRepository struct {
	db *gorm.DB
}

func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{db: db}
}

func (r *RequestLogRepository) Create(ctx context.Context, log *domain.RequestLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *RequestLogRepository) GetRequestsPerMinute(ctx context.Context, last24h bool) ([]domain.MetricResponse, error) {
	var results []domain.MetricResponse

	query := r.db.WithContext(ctx).Model(&domain.RequestLog{}).
		Select("date_trunc('minute', created_at) as timestamp, count(*) as count").
		Group("timestamp").
		Order("timestamp DESC")

	if last24h {
		query = query.Where("created_at > ?", time.Now().Add(-24*time.Hour))
	}

	err := query.Find(&results).Error
	return results, err
}
