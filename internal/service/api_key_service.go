package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) CreateAPIKey(ctx context.Context, req domain.APIKeyCreateRequest) (*domain.APIKey, error) {
	// Generate a random key
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	plainKey := hex.EncodeToString(bytes)

	// Hash the key using SHA-256
	hash := sha256.Sum256([]byte(plainKey))
	hashedKey := hex.EncodeToString(hash[:])

	prefix := fmt.Sprintf("env_%s...", plainKey[:8])

	apiKey := &domain.APIKey{
		Name:      req.Name,
		HashedKey: hashedKey,
		PlainKey:  plainKey, // This will be returned to the user only once
		Prefix:    prefix,
		Status:    "active",
		ExpiresAt: req.ExpiresAt,
	}

	// Add applications if provided
	if len(req.ApplicationIDs) > 0 {
		for _, id := range req.ApplicationIDs {
			apiKey.Applications = append(apiKey.Applications, domain.Application{ID: id})
		}
	}

	if err := s.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *APIKeyService) GetAllAPIKeys(ctx context.Context) ([]domain.APIKey, error) {
	return s.repo.GetAll(ctx)
}

func (s *APIKeyService) RevokeAPIKey(ctx context.Context, id uint) error {
	return s.repo.UpdateStatus(ctx, id, "revoked")
}

func (s *APIKeyService) ValidateKey(ctx context.Context, plainKey string) (*domain.APIKey, error) {
	hash := sha256.Sum256([]byte(plainKey))
	hashedKey := hex.EncodeToString(hash[:])

	apiKey, err := s.repo.FindByHashedKey(ctx, hashedKey)
	if err != nil {
		return nil, err
	}

	// Check expiration
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("API Key expired")
	}

	// Update last used
	_ = s.repo.UpdateLastUsed(ctx, apiKey.ID)

	return apiKey, nil
}
