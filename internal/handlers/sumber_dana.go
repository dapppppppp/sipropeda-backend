package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/master"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type SumberDanaHandler struct {
	service master.SumberDanaService
}

func ProvideSumberDanaHandler(service master.SumberDanaService) SumberDanaHandler {
	return SumberDanaHandler{service: service}
}

func (h *SumberDanaHandler) Router(r chi.Router) {
	r.Route("/sumber-dana", func(rc chi.Router) {
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

// Create menambah data Sumber Dana
// @Summary Tambah data Sumber Dana
// @Tags Sumber Dana
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body master.RequestSumberDana true "Data Sumber Dana"
// @Success 201 {object} response.Base
// @Router /v1/sumber-dana [post]
func (h *SumberDanaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestSumberDana
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.service.Create(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, req)
}

// ResolveAll mengambil semua data Sumber Dana
// @Summary Ambil semua data Sumber Dana
// @Tags Sumber Dana
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/sumber-dana [get]
func (h *SumberDanaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID mengambil data Sumber Dana berdasarkan ID
// @Summary Ambil detail Sumber Dana by ID
// @Tags Sumber Dana
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Sumber Dana"
// @Success 200 {object} response.Base
// @Router /v1/sumber-dana/{id} [get]
func (h *SumberDanaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
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

// Update mengubah data Sumber Dana
// @Summary Update data Sumber Dana
// @Tags Sumber Dana
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Sumber Dana"
// @Param body body master.RequestSumberDana true "Data yang akan diedit"
// @Success 200 {object} response.Base
// @Router /v1/sumber-dana/{id} [put]
func (h *SumberDanaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req master.RequestSumberDana
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.service.Update(id, req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Sumber Dana successfully updated"})
}

// DeleteSoft menghapus data Sumber Dana
// @Summary Hapus data Sumber Dana (Soft Delete)
// @Tags Sumber Dana
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Sumber Dana"
// @Success 200 {object} response.Base
// @Router /v1/sumber-dana/{id} [delete]
func (h *SumberDanaHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(id); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Sumber Dana successfully deleted"})
}