package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/repository"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) CreateAPIKey(ctx context.Context, name string) (*domain.APIKey, error) {
	// Generate a random key
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	key := hex.EncodeToString(bytes)
	prefix := fmt.Sprintf("pk_%s...", key[:8])

	apiKey := &domain.APIKey{
		Name:   name,
		Key:    key,
		Prefix: prefix,
		Status: "active",
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
