package domain

import (
	"time"

	"gorm.io/gorm"
)

// --- Role Types ---

type RoleType string

const (
	RoleAdmin  RoleType = "admin"
	RoleEditor RoleType = "editor"
	RoleViewer RoleType = "viewer"
)

// --- User ---

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name" validate:"required,min=2"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      RoleType       `gorm:"type:varchar(20);not null;default:'viewer'" json:"role" validate:"required,oneof=admin editor viewer"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Application ---

type Application struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name" validate:"required,min=2"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Environment ---

type Environment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name" validate:"required,min=2"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Configuration ---

type Configuration struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Key           string         `gorm:"not null;index:idx_config_lookup" json:"key" validate:"required"`
	Value         string         `gorm:"type:text;not null" json:"value" validate:"required"`
	ApplicationID uint           `gorm:"not null;index:idx_config_lookup" json:"application_id" validate:"required"`
	EnvironmentID uint           `gorm:"not null;index:idx_config_lookup" json:"environment_id" validate:"required"`
	Version       int            `gorm:"not null;default:1" json:"version"`
	Status        string         `gorm:"not null;default:'active'" json:"status"` // active, archived
	CreatedBy     string         `json:"created_by_name"`
	Description   string         `gorm:"type:text" json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Application Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Environment Environment `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
}

// --- Config Version (audit history) ---

type ConfigVersion struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ConfigurationID uint      `gorm:"not null;index" json:"configuration_id"`
	Key             string    `gorm:"not null" json:"key"`
	Value           string    `gorm:"type:text" json:"value"`
	Version         int       `gorm:"not null" json:"version"`
	ChangedByUserID uint      `json:"changed_by_user_id"`
	ChangedByName   string    `json:"changed_by_name"`
	CreatedAt       time.Time `json:"created_at"`
}

// --- Audit Log ---

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

// --- API Key ---

type APIKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Key       string         `gorm:"uniqueIndex;not null" json:"key"`
	Prefix    string         `json:"prefix"`
	Status    string         `json:"status"` // active, revoked
	LastUsed  *time.Time     `json:"last_used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Dashboard Stats ---

type DashboardStats struct {
	TotalApps         int64      `json:"total_apps"`
	TotalConfigs      int64      `json:"total_configs"`
	ActiveConfigs     int64      `json:"active_configs"`
	TotalEnvironments int64      `json:"total_environments"`
	RecentUpdates     []AuditLog `json:"recent_updates"`
}

// --- DTOs ---

type RegisterRequest struct {
	Name     string   `json:"name" validate:"required,min=2"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     RoleType `json:"role" validate:"required,oneof=admin editor viewer"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateUserRequest struct {
	Name  string   `json:"name" validate:"omitempty,min=2"`
	Email string   `json:"email" validate:"omitempty,email"`
	Role  RoleType `json:"role" validate:"omitempty,oneof=admin editor viewer"`
}

type UpdateConfigRequest struct {
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}
