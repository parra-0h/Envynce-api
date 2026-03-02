package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/service"
	"github.com/hans/config-service/pkg/utils"
)

var validate = validator.New()

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authSvc}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := validate.Struct(req); err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, user, "User registered successfully")
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := validate.Struct(req); err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp, "Login successful")
}
