package repository

import (
	"github.com/hans/config-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
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

	// Seed Data
	SeedData(db)

	return db, nil
}

func SeedData(db *gorm.DB) {
	// Check if user exists
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "parrahans70@gmail.com").Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&domain.User{
			Name:     "Hans Parra",
			Email:    "parrahans70@gmail.com",
			Password: string(hashed),
			Role:     domain.RoleAdmin,
		})
	}
}
