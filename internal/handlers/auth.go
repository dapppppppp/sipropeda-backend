package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type AuthHandler struct {
	service auth.UserService
}

func ProvideAuthHandler(service auth.UserService) AuthHandler {
	return AuthHandler{service: service}
}

// Router untuk Endpoint Public (Login)
func (h *AuthHandler) Router(r chi.Router) {
	r.Post("/login", h.Login)
}

// UserRouter untuk Endpoint Private CRUD User
func (h *AuthHandler) UserRouter(r chi.Router) {
	r.Route("/user", func(rc chi.Router) {
		rc.Get("/", h.ResolveAll)
		rc.Post("/", h.Create)
		rc.Put("/{id}", h.Update)
		rc.Get("/{id}", h.ResolveByID)
		rc.Delete("/{id}", h.DeleteSoft)
	})
}

// Login User
// @Summary Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.LoginRequest true "Kredensial Login"
// @Success 200 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /v1/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, err)
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	response.WithJSON(w, http.StatusOK, resp)
}

// Create User
// @Summary Tambah User Baru
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestUserFormat true "Data User"
// @Success 201 {object} response.Base
// @Router /v1/user [post]
func (h *AuthHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req auth.RequestUserFormat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body: "+err.Error()))
		return
	}

	if err := h.service.Create(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, req)
}

// ResolveAll User
// @Summary Ambil Semua Data User
// @Tags User
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/user [get]
func (h *AuthHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID User
// @Summary Ambil detail User by ID
// @Tags User
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID User"
// @Success 200 {object} response.Base
// @Router /v1/user/{id} [get]
func (h *AuthHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	parsedID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, errors.New("invalid UUID format"))
		return
	}

	data, err := h.service.ResolveByID(parsedID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// Update User
// @Summary Perbarui data User
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID User"
// @Param body body auth.RequestUserFormat true "Data User yang diperbarui (kosongkan password jika tidak ingin ganti)"
// @Success 200 {object} response.Base
// @Router /v1/user/{id} [put]
func (h *AuthHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req auth.RequestUserFormat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body: "+err.Error()))
		return
	}

	if err := h.service.Update(id, req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "User successfully updated"})
}

// DeleteSoft User
// @Summary Hapus User by ID (Soft Delete)
// @Tags User
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID User"
// @Success 200 {object} response.Base
// @Router /v1/user/{id} [delete]
func (h *AuthHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(id); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "User successfully deleted"})
}