package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/middleware"
	"github.com/hans/config-service/internal/service"
	"github.com/hans/config-service/pkg/utils"
)

type BaseHandler struct {
	appService    *service.ApplicationService
	envService    *service.EnvironmentService
	configService *service.ConfigurationService
	apiKeyService *service.APIKeyService
}

func NewBaseHandler(appSvc *service.ApplicationService, envSvc *service.EnvironmentService, configSvc *service.ConfigurationService, apiKeySvc *service.APIKeyService) *BaseHandler {
	return &BaseHandler{
		appService:    appSvc,
		envService:    envSvc,
		configService: configSvc,
		apiKeyService: apiKeySvc,
	}
}

// HealthCheck GET /health
func (h *BaseHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"}, "Service is healthy")
}

// ---------------- Applications ----------------

// POST /api/v1/applications
func (h *BaseHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var app domain.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if app.Name == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}
	if err := h.appService.CreateApplication(r.Context(), &app); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusCreated, app, "Application created successfully")
}

// GET /api/v1/applications
func (h *BaseHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appService.GetAllApplications(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, apps, "")
}

// GET /api/v1/applications/{id}
func (h *BaseHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}
	app, err := h.appService.GetApplicationByID(r.Context(), uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, app, "")
}

// PUT /api/v1/applications/{id}
func (h *BaseHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	app, err := h.appService.UpdateApplication(r.Context(), uint(id), body.Name, body.Description)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, app, "Application updated successfully")
}

// DELETE /api/v1/applications/{id}
func (h *BaseHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid application ID")
		return
	}
	if err := h.appService.DeleteApplication(r.Context(), uint(id)); err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil, "Application deleted successfully")
}

// ---------------- Environments ----------------

// POST /api/v1/environments
func (h *BaseHandler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var env domain.Environment
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if env.Name == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}
	if err := h.envService.CreateEnvironment(r.Context(), &env); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusCreated, env, "Environment created successfully")
}

// GET /api/v1/environments
func (h *BaseHandler) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := h.envService.GetAllEnvironments(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, envs, "")
}

// GET /api/v1/environments/{id}
func (h *BaseHandler) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}
	env, err := h.envService.GetEnvironmentByID(r.Context(), uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, env, "")
}

// PUT /api/v1/environments/{id}
func (h *BaseHandler) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	env, err := h.envService.UpdateEnvironment(r.Context(), uint(id), body.Name, body.Description)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, env, "Environment updated successfully")
}

// DELETE /api/v1/environments/{id}
func (h *BaseHandler) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid environment ID")
		return
	}
	if err := h.envService.DeleteEnvironment(r.Context(), uint(id)); err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil, "Environment deleted successfully")
}

// ---------------- Configurations ----------------

// POST /api/v1/configs
func (h *BaseHandler) CreateConfiguration(w http.ResponseWriter, r *http.Request) {
	var config domain.Configuration
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Printf("BACKEND: CreateConfiguration decode error: %v", err)
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if config.Key == "" || config.Value == "" || config.ApplicationID == 0 || config.EnvironmentID == 0 {
		utils.JSONError(w, http.StatusUnprocessableEntity, "key, value, application_id and environment_id are required")
		return
	}

	userID := middleware.GetUserID(r.Context())
	userName := middleware.GetUserName(r.Context())

	if err := h.configService.CreateConfiguration(r.Context(), &config, userID, userName); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusCreated, config, "Configuration created successfully")
}

// GET /api/v1/configs
func (h *BaseHandler) ListConfigurations(w http.ResponseWriter, r *http.Request) {
	appID, _ := strconv.ParseUint(r.URL.Query().Get("application_id"), 10, 64)
	envID, _ := strconv.ParseUint(r.URL.Query().Get("environment_id"), 10, 64)
	search := r.URL.Query().Get("search")

	log.Printf("BACKEND: ListConfigurations called for appID=%d, envID=%d, search=%s", appID, envID, search)

	if appID == 0 || envID == 0 {
		utils.JSONError(w, http.StatusBadRequest, "application_id and environment_id are required")
		return
	}

	var configs []domain.Configuration
	var err error
	if search != "" {
		configs, err = h.configService.SearchConfigurations(r.Context(), uint(appID), uint(envID), search)
	} else {
		configs, err = h.configService.GetActiveConfigs(r.Context(), uint(appID), uint(envID))
	}

	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, configs, "")
}

// GET /api/v1/configs/{id}
func (h *BaseHandler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}
	config, err := h.configService.GetByID(r.Context(), uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, config, "")
}

// PUT /api/v1/configs/{id}
func (h *BaseHandler) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}
	var req domain.UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if req.Value == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, "value is required")
		return
	}

	userID := middleware.GetUserID(r.Context())
	userName := middleware.GetUserName(r.Context())

	config, err := h.configService.UpdateConfiguration(r.Context(), uint(id), &req, userID, userName)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, config, "Configuration updated successfully")
}

// DELETE /api/v1/configs/{id}
func (h *BaseHandler) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}
	if err := h.configService.DeleteConfiguration(r.Context(), uint(id)); err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil, "Configuration deleted successfully")
}

// GET /api/v1/configs/{id}/versions
func (h *BaseHandler) GetConfigVersions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}
	versions, err := h.configService.GetConfigVersions(r.Context(), uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, versions, "")
}

func (h *BaseHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.configService.GetAuditLogs(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, logs, "")
}

func (h *BaseHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.configService.GetDashboardStats(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, stats, "")
}

func (h *BaseHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Simple check for demonstration
	if payload.Email == "parrahans70@gmail.com" {
		// In a real app we'd check DB and hash passwords
		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Et9HFBa_nSYUCPUoKndp1vR-fS-m5W4vI-b6-M_b4-I",
			"user": map[string]interface{}{
				"id":    "1",
				"email": payload.Email,
				"name":  "Hans Parra",
				"role":  "admin",
			},
		}, "Login successful")
		return
	}

	utils.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
}

func (h *BaseHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	var payload domain.APIKeyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if payload.Name == "" {
		utils.JSONError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}

	key, err := h.apiKeyService.CreateAPIKey(r.Context(), payload)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusCreated, key, "API Key generated successfully")
}

func (h *BaseHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := h.apiKeyService.GetAllAPIKeys(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, keys, "")
}

func (h *BaseHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if id == 0 {
		utils.JSONError(w, http.StatusBadRequest, "Invalid API Key ID")
		return
	}

	if err := h.apiKeyService.RevokeAPIKey(r.Context(), uint(id)); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil, "API Key revoked successfully")
}

func (h *BaseHandler) RestoreConfigVersion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if id == 0 {
		utils.JSONError(w, http.StatusBadRequest, "Invalid version ID")
		return
	}

	userID := middleware.GetUserID(r.Context())
	userName := middleware.GetUserName(r.Context())

	config, err := h.configService.RestoreVersion(r.Context(), uint(id), userID, userName)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, config, "Configuration restored successfully")
}
