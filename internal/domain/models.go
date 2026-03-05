package domain

import (
	"time"

	"gorm.io/gorm"
)

// --- Role Types ---

type RoleType string

const (
	RoleAdmin     RoleType = "admin"
	RoleDeveloper RoleType = "developer"
	RoleViewer    RoleType = "viewer"
)

// --- User ---

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name" validate:"required,min=2"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      RoleType       `gorm:"type:varchar(20);not null;default:'viewer'" json:"role" validate:"required,oneof=admin developer viewer"`
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
	ID              uint           `gorm:"primaryKey" json:"id"`
	ConfigurationID uint           `gorm:"not null;index" json:"configuration_id"`
	Key             string         `gorm:"not null" json:"key"`
	Value           string         `gorm:"type:text" json:"value"`
	Description     string         `gorm:"type:text" json:"description"`
	Active          bool           `json:"active"`
	VersionNumber   int            `gorm:"not null;default:0" json:"version_number"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Audit Log ---

type AuditLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	Action     string         `json:"action"` // CREATE, UPDATE, DELETE, LOGIN, LOGIN_FAILED, GENERATE_KEY, REVOKE_KEY
	EntityType string         `json:"entity_type"`
	EntityID   uint           `json:"entity_id"`
	Metadata   string         `gorm:"type:jsonb" json:"metadata"` // JSON string for flexibility
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- API Key ---

type APIKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	HashedKey string         `gorm:"uniqueIndex;not null;default:''" json:"-"`
	PlainKey  string         `gorm:"-" json:"key,omitempty"` // Only for response ONCE
	Prefix    string         `gorm:"default:''" json:"prefix"`
	Status    string         `gorm:"default:'active'" json:"status"` // active, revoked
	ExpiresAt *time.Time     `json:"expires_at"`
	LastUsed  *time.Time     `json:"last_used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Applications []Application `gorm:"many2many:api_key_applications;" json:"applications"`
}

// --- Request Log ---

type RequestLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	APIKeyID      uint      `gorm:"index" json:"api_key_id"`
	ApplicationID uint      `gorm:"index" json:"application_id"`
	EnvironmentID uint      `gorm:"index" json:"environment_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// --- DTOs ---

type DashboardStats struct {
	TotalApps         int64      `json:"total_apps"`
	TotalConfigs      int64      `json:"total_configs"`
	ActiveConfigs     int64      `json:"active_configs"`
	TotalEnvironments int64      `json:"total_environments"`
	RecentUpdates     []AuditLog `json:"recent_updates"`
}

type RegisterRequest struct {
	Name     string   `json:"name" validate:"required,min=2"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     RoleType `json:"role" validate:"required,oneof=admin developer viewer"`
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
	Role  RoleType `json:"role" validate:"omitempty,oneof=admin developer viewer"`
}

type UpdateConfigRequest struct {
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}

type APIKeyCreateRequest struct {
	Name           string     `json:"name" validate:"required"`
	ApplicationIDs []uint     `json:"application_ids"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

type MetricResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}
