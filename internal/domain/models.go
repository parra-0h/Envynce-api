package domain

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex:idx_applications_name" json:"name" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Environment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex:uni_environments_name;not null" json:"name" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Configuration struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Key           string         `gorm:"not null;index:idx_config_key" json:"key" validate:"required"`
	Value         string         `gorm:"type:text" json:"value" validate:"required"`
	ApplicationID uint           `gorm:"not null" json:"application_id" validate:"required"`
	EnvironmentID uint           `gorm:"not null" json:"environment_id" validate:"required"`
	Version       int            `gorm:"not null;default:1" json:"version"`
	Status        string         `gorm:"not null;default:'active'" json:"status"` // active, archived
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Application Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Environment Environment `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Action    string    `json:"action"` // CREATE, UPDATE, DELETE
	Entity    string    `json:"entity"` // Application, Environment, Configuration
	EntityID  uint      `json:"entity_id"`
	OldValue  string    `gorm:"type:text" json:"old_value"`
	NewValue  string    `gorm:"type:text" json:"new_value"`
	ChangedBy string    `json:"changed_by"` // API Key or User identifier
	CreatedAt time.Time `json:"created_at"`
}
