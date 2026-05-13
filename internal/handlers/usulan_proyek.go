package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/transaction"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type UsulanProyekHandler struct {
	service transaction.UsulanProyekService
}

func ProvideUsulanProyekHandler(service transaction.UsulanProyekService) UsulanProyekHandler {
	return UsulanProyekHandler{service: service}
}

func (h *UsulanProyekHandler) Router(r chi.Router) {
	r.Route("/usulan-proyek", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/", h.ResolveAll)
			protected.Post("/", h.Create)
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/{id}", h.Update)
			protected.Delete("/{id}", h.DeleteSoft)
		})
	})
}

// Create menambah data Usulan Proyek
// @Summary Tambah data Usulan Proyek
// @Tags Usulan Proyek
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body transaction.RequestUsulanProyek true "Data Usulan Proyek"
// @Success 201 {object} response.Base
// @Router /v1/usulan-proyek [post]
func (h *UsulanProyekHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req transaction.RequestUsulanProyek
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}
	req.UserID = userID

	if err := h.service.Create(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, req)
}

// ResolveAll mengambil semua data Usulan Proyek
// @Summary Ambil semua data Usulan Proyek
// @Tags Usulan Proyek
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek [get]
func (h *UsulanProyekHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID mengambil data Usulan Proyek berdasarkan ID
// @Summary Ambil detail Usulan Proyek by ID
// @Tags Usulan Proyek
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Usulan Proyek"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/{id} [get]
func (h *UsulanProyekHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
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

// Update mengubah data Usulan Proyek
// @Summary Update data Usulan Proyek
// @Tags Usulan Proyek
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Usulan Proyek"
// @Param body body transaction.RequestUsulanProyek true "Data yang akan diedit"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/{id} [put]
func (h *UsulanProyekHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req transaction.RequestUsulanProyek
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}
	req.UserID = userID

	if err := h.service.Update(id, req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Usulan Proyek successfully updated"})
}

// DeleteSoft menghapus data Usulan Proyek
// @Summary Hapus data Usulan Proyek (Soft Delete)
// @Tags Usulan Proyek
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Usulan Proyek"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/{id} [delete]
func (h *UsulanProyekHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Usulan Proyek successfully deleted"})
}