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

type KriteriaHandler struct {
	service master.KriteriaService
}

// ProvideKriteriaHandler adalah constructor untuk injeksi dependensi
func ProvideKriteriaHandler(service master.KriteriaService) KriteriaHandler {
	return KriteriaHandler{service: service}
}

// Router mendaftarkan rute-rute khusus untuk Kriteria
func (h *KriteriaHandler) Router(r chi.Router) {
	r.Route("/kriteria", func(rc chi.Router) {
		rc.Get("/", h.ResolveAll)
		rc.Post("/", h.Create)
		rc.Put("/{id}", h.Update)
		rc.Get("/{id}", h.ResolveByID)
		rc.Delete("/{id}", h.DeleteSoft)
	})
}

// Create menambahkan data Kriteria baru.
// @Summary Tambah data Kriteria baru
// @Description Endpoint ini digunakan untuk menambahkan Kriteria baru ke database.
// @Tags Kriteria
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body master.RequestKriteriaFormat true "Data Kriteria"
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Router /v1/kriteria [post]
func (h *KriteriaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestKriteriaFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, errors.New("invalid JSON body: "+err.Error()))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}

	req.UserID = userID
	err = h.service.Create(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, req)
}

// ResolveAll mengambil daftar Kriteria
// @Summary Ambil Semua Data Kriteria
// @Description Endpoint untuk mengambil daftar kriteria yang aktif
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/kriteria [get]
func (h *KriteriaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// Update memperbarui data Kriteria berdasarkan ID.
// @Summary Perbarui data Kriteria
// @Description Endpoint ini digunakan untuk memperbarui data Kriteria berdasarkan ID yang dikirimkan di path.
// @Tags Kriteria
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Param body body master.RequestKriteriaFormat true "Data Kriteria yang diperbarui"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/kriteria/{id} [put]
func (h *KriteriaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WithError(w, errors.New("missing id in path"))
		return
	}

	var req master.RequestKriteriaFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, errors.New("invalid JSON body: "+err.Error()))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}

	req.UserID = userID
	err = h.service.Update(id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Kriteria successfully updated"})
}

// ResolveByID mengambil data Kriteria berdasarkan ID.
// @Summary Ambil detail Kriteria by ID
// @Description Endpoint ini digunakan untuk mengambil satu data Kriteria berdasarkan ID.
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/kriteria/{id} [get]
func (h *KriteriaHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
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

// DeleteSoft menghapus data Kriteria secara soft delete berdasarkan ID.
// @Summary Hapus Kriteria by ID (Soft Delete)
// @Description Endpoint ini digunakan untuk menghapus data Kriteria berdasarkan ID (soft delete).
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/kriteria/{id} [delete]
func (h *KriteriaHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WithError(w, errors.New("missing id in path"))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak ditemukan"})
		return
	}

	err := h.service.Delete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Kriteria successfully deleted"})
}