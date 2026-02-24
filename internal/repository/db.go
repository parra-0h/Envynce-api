package repository

import (
	"github.com/hans/config-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&domain.Application{},
		&domain.Environment{},
		&domain.Configuration{},
		&domain.AuditLog{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
