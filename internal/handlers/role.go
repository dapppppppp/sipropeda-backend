package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type RoleHandler struct {
	RoleService auth.RoleService
}

func ProvideRoleHandler(roleService auth.RoleService) RoleHandler {
	return RoleHandler{RoleService: roleService}
}

func (h *RoleHandler) Router(r chi.Router) {
	r.Route("/roles", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected) // Middleware dari Sipropeda
			protected.Get("/", h.ResolveAll)
			protected.Get("/all", h.GetAllData) // Sama dengan ResolveAll untuk saat ini
			protected.Post("/", h.CreateRole)
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/{id}", h.UpdateRole)
			protected.Delete("/{id}", h.DeleteRole)
		})
	})
}

// ResolveAll mengambil semua data Role
// @Summary Ambil semua data Role
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/roles [get]
func (h *RoleHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	// (Pagination dilewati agar kompatibel dengan modul saat ini)
	data, err := h.RoleService.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// GetAllData mengambil semua data Role
// @Summary Ambil semua data Role
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/roles/all [get]
func (h *RoleHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := h.RoleService.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// CreateRole menambah data Role
// @Summary Tambah data Role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestRole true "Role yang akan ditambahkan"
// @Success 201 {object} response.Base
// @Router /v1/roles [post]
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req auth.RequestRole
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.RoleService.CreateRole(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, req)
}

// UpdateRole mengubah data Role
// @Summary Update data Role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Role"
// @Param body body auth.RequestRole true "Role yang akan diedit"
// @Success 200 {object} response.Base
// @Router /v1/roles/{id} [put]
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req auth.RequestRole
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.RoleService.UpdateRole(id, req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Role successfully updated"})
}

// ResolveByID mengambil data Role berdasarkan ID
// @Summary Ambil detail Role by ID
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Role"
// @Success 200 {object} response.Base
// @Router /v1/roles/{id} [get]
func (h *RoleHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	parsedID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, errors.New("invalid UUID format"))
		return
	}

	data, err := h.RoleService.ResolveByID(parsedID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// DeleteRole menghapus data Role
// @Summary Hapus data Role (Soft Delete)
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Role"
// @Success 200 {object} response.Base
// @Router /v1/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.RoleService.DeleteRole(id); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Role successfully deleted"})
}