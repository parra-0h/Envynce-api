package repository

import (
	"context"
	"encoding/json"

	"github.com/hans/config-service/internal/domain"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepository) GetAll(ctx context.Context, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) LogAction(ctx context.Context, userID uint, action, entityType string, entityID uint, metadata interface{}) error {
	metaJSON, _ := json.Marshal(metadata)
	log := &domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Metadata:   string(metaJSON),
	}
	return r.Create(ctx, log)
}
