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

type PaguAnggaranHandler struct {
	service transaction.PaguAnggaranService
}

func ProvidePaguAnggaranHandler(service transaction.PaguAnggaranService) PaguAnggaranHandler {
	return PaguAnggaranHandler{service: service}
}

func (h *PaguAnggaranHandler) Router(r chi.Router) {
	r.Route("/pagu-anggaran", func(rc chi.Router) {
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

// Create menambah data Pagu Anggaran
// @Summary Tambah data Pagu Anggaran
// @Tags Pagu Anggaran
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body transaction.RequestPaguAnggaran true "Data Pagu Anggaran"
// @Success 201 {object} response.Base
// @Router /v1/pagu-anggaran [post]
func (h *PaguAnggaranHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req transaction.RequestPaguAnggaran
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

// ResolveAll mengambil semua data Pagu Anggaran
// @Summary Ambil semua data Pagu Anggaran
// @Tags Pagu Anggaran
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/pagu-anggaran [get]
func (h *PaguAnggaranHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID mengambil data Pagu Anggaran berdasarkan ID
// @Summary Ambil detail Pagu Anggaran by ID
// @Tags Pagu Anggaran
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Pagu Anggaran"
// @Success 200 {object} response.Base
// @Router /v1/pagu-anggaran/{id} [get]
func (h *PaguAnggaranHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
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

// Update mengubah data Pagu Anggaran
// @Summary Update data Pagu Anggaran
// @Tags Pagu Anggaran
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Pagu Anggaran"
// @Param body body transaction.RequestPaguAnggaran true "Data yang akan diedit"
// @Success 200 {object} response.Base
// @Router /v1/pagu-anggaran/{id} [put]
func (h *PaguAnggaranHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req transaction.RequestPaguAnggaran
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
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Pagu Anggaran successfully updated"})
}

// DeleteSoft menghapus data Pagu Anggaran
// @Summary Hapus data Pagu Anggaran (Soft Delete)
// @Tags Pagu Anggaran
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Pagu Anggaran"
// @Success 200 {object} response.Base
// @Router /v1/pagu-anggaran/{id} [delete]
func (h *PaguAnggaranHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
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
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Pagu Anggaran successfully deleted"})
}