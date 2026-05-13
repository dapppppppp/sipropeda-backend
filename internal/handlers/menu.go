package handlers

import (
	"encoding/json"
	"net/http"

	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
)

type MenuHandler struct {
	service auth.MenuService
}

func ProvideMenuHandler(service auth.MenuService) MenuHandler {
	return MenuHandler{service: service}
}

func (h *MenuHandler) Router(r chi.Router) {
	r.Route("/menu", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/", h.ResolveAll)
			protected.Post("/", h.Create)
			protected.Put("/{id}", h.Update)
			protected.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/menu-role", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/", h.ResolveMenuByRoleID)
			protected.Post("/bulk", h.SaveBulkMenuRole)
		})
	})
}

// ResolveAll Menu
// @Summary Ambil semua Master Menu
// @Tags Menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/menu [get]
func (h *MenuHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAllMenu()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// Create Menu
// @Summary Tambah Master Menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestMenuFormat true "Data Menu"
// @Success 201 {object} response.Base
// @Router /v1/menu [post]
func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req auth.RequestMenuFormat
	json.NewDecoder(r.Body).Decode(&req)
	err := h.service.CreateMenu(req)
	if err != nil { response.WithError(w, err); return }
	response.WithJSON(w, http.StatusCreated, req)
}

// Update Menu
// @Summary Update Master Menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Menu"
// @Param body body auth.RequestMenuFormat true "Data Menu"
// @Success 200 {object} response.Base
// @Router /v1/menu/{id} [put]
func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req auth.RequestMenuFormat
	json.NewDecoder(r.Body).Decode(&req)
	err := h.service.UpdateMenu(id, req)
	if err != nil { response.WithError(w, err); return }
	response.WithJSON(w, http.StatusOK, "success")
}

// Delete Menu
// @Summary Hapus Master Menu
// @Tags Menus
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Menu"
// @Success 200 {object} response.Base
// @Router /v1/menu/{id} [delete]
func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteMenu(id)
	if err != nil { response.WithError(w, err); return }
	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveMenuByRoleID
// @Summary Ambil Hierarki Menu by Role ID
// @Tags Menus Role
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param roleId query string true "Role ID"
// @Success 200 {object} response.Base
// @Router /v1/menu-role [get]
func (h *MenuHandler) ResolveMenuByRoleID(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	data, err := h.service.ResolveMenuByRoleID(roleID)
	if err != nil { response.WithError(w, err); return }
	response.WithJSON(w, http.StatusOK, data)
}

// SaveBulkMenuRole
// @Summary Pasang Akses Menu ke Role
// @Tags Menus Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body auth.RequestBulkMenuRole true "Data Mapping"
// @Success 200 {object} response.Base
// @Router /v1/menu-role/bulk [post]
func (h *MenuHandler) SaveBulkMenuRole(w http.ResponseWriter, r *http.Request) {
	var req auth.RequestBulkMenuRole
	json.NewDecoder(r.Body).Decode(&req)
	err := h.service.SaveBulkMenuRole(req)
	if err != nil { response.WithError(w, err); return }
	response.WithJSON(w, http.StatusOK, "success")
}