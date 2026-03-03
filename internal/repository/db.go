package repository

import (
	"github.com/hans/config-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Application{},
		&domain.Environment{},
		&domain.Configuration{},
		&domain.ConfigVersion{},
		&domain.AuditLog{},
		&domain.APIKey{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
