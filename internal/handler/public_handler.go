package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hans/config-service/internal/middleware"
	"github.com/hans/config-service/internal/service"
	"github.com/hans/config-service/pkg/utils"
)

type PublicHandler struct {
	configSvc  *service.ConfigurationService
	metricsSvc *service.MetricsService
}

func NewPublicHandler(configSvc *service.ConfigurationService, metricsSvc *service.MetricsService) *PublicHandler {
	return &PublicHandler{
		configSvc:  configSvc,
		metricsSvc: metricsSvc,
	}
}

// GetConfig GET /api/v1/config
// @Summary Get configuration as flat JSON
// @Description Returns active configurations for a specific application and environment.
// @Tags Public
// @Security BearerAuth
// @Param application query string true "Application Name"
// @Param environment query string true "Environment Name"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /config [get]
func (h *PublicHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("application")
	env := r.URL.Query().Get("environment")

	if app == "" || env == "" {
		utils.JSONError(w, http.StatusBadRequest, "application and environment are required query parameters")
		return
	}

	configs, appID, envID, err := h.configSvc.GetConfigsAsMap(r.Context(), app, env)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	// Scoping check (Level 3)
	apiKeyID := middleware.GetAPIKeyID(r.Context())
	apiKeyApps := middleware.GetAPIKeyApps(r.Context())

	if len(apiKeyApps) > 0 {
		found := false
		for _, a := range apiKeyApps {
			if a.ID == appID {
				found = true
				break
			}
		}
		if !found {
			utils.JSONError(w, http.StatusForbidden, "API Key does not have access to this application")
			return
		}
	}

	// Log metrics (Level 1)
	go func() {
		_ = h.metricsSvc.LogRequest(context.Background(), apiKeyID, appID, envID)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}
