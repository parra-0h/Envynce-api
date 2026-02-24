package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/service"
	"github.com/hans/config-service/pkg/utils"
)

type BaseHandler struct {
	appService    *service.ApplicationService
	envService    *service.EnvironmentService
	configService *service.ConfigurationService
}

func NewBaseHandler(appSvc *service.ApplicationService, envSvc *service.EnvironmentService, configSvc *service.ConfigurationService) *BaseHandler {
	return &BaseHandler{
		appService:    appSvc,
		envService:    envSvc,
		configService: configSvc,
	}
}

// Health Check
func (h *BaseHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, nil, "Service is healthy")
}

// Applications
func (h *BaseHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var app domain.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.appService.CreateApplication(r.Context(), &app); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, app, "Application created successfully")
}

func (h *BaseHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appService.GetAllApplications(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, apps, "")
}

// Environments
func (h *BaseHandler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var env domain.Environment
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.envService.CreateEnvironment(r.Context(), &env); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, env, "Environment created successfully")
}

func (h *BaseHandler) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := h.envService.GetAllEnvironments(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, envs, "")
}

// Configurations
func (h *BaseHandler) CreateConfiguration(w http.ResponseWriter, r *http.Request) {
	var config domain.Configuration
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.configService.CreateConfiguration(r.Context(), &config); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, config, "Configuration created successfully")
}

func (h *BaseHandler) ListConfigurations(w http.ResponseWriter, r *http.Request) {
	appID, _ := strconv.Atoi(r.URL.Query().Get("application_id"))
	envID, _ := strconv.Atoi(r.URL.Query().Get("environment_id"))

	if appID == 0 || envID == 0 {
		utils.JSONError(w, http.StatusBadRequest, "application_id and environment_id are required")
		return
	}

	configs, err := h.configService.GetActiveConfigs(r.Context(), uint(appID), uint(envID))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, configs, "")
}

func (h *BaseHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.configService.GetAuditLogs(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, logs, "")
}
