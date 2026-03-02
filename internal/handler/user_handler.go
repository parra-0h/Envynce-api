package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/service"
	"github.com/hans/config-service/pkg/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userService: userSvc}
}

// GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, users, "")
}

// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, err := h.userService.GetUserByID(r.Context(), uint(id))
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, user, "")
}

// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := validate.Struct(req); err != nil {
		utils.JSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	user, err := h.userService.UpdateUser(r.Context(), uint(id), &req)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, user, "User updated successfully")
}

// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := h.userService.DeleteUser(r.Context(), uint(id)); err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil, "User deleted successfully")
}
