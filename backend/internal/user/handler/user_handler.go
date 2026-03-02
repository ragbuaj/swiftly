package handler

import (
	"encoding/json"
	"net/http"
	"swiftly/backend/internal/pkg/response"
	"swiftly/backend/internal/user"
	"swiftly/backend/internal/user/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users/profile", h.GetUser)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	u, err := h.service.CreateUser(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusCreated, "User created successfully", u)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	u, err := h.service.GetUserByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found", nil)
		return
	}

	response.Success(w, http.StatusOK, "User retrieved successfully", u)
}
