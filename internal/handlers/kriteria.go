package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"sipropeda-backend/internal/domain/master"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type KriteriaHandler struct {
	service master.KriteriaService
}

func ProvideKriteriaHandler(service master.KriteriaService) KriteriaHandler {
	return KriteriaHandler{service: service}
}

func (h *KriteriaHandler) Router(r chi.Router) {
	r.Route("/kriteria", func(rc chi.Router) {
		// Pasang JWT Middleware untuk semua endpoint kriteria
		rc.Use(middleware.JWTProtected)
		rc.Get("/", h.ResolveAll)
		rc.Get("/page", h.ResolvePaging) // Endpoint baru pagination (Taruh di atas /{id})
		rc.Post("/", h.Create)
		rc.Get("/{id}", h.ResolveByID)
		rc.Put("/{id}", h.Update)
		rc.Delete("/{id}", h.DeleteSoft)
	})
}

// Create menambahkan data Kriteria baru.
// @Summary Tambah data Kriteria baru
// @Tags Kriteria
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body master.RequestKriteriaFormat true "Data Kriteria"
// @Success 201 {object} response.Base
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

// ResolveAll mengambil daftar Kriteria (Semua data, untuk kalkulasi/dropdown)
// @Summary Ambil Semua Data Kriteria
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/kriteria [get]
func (h *KriteriaHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolvePaging mengambil data Kriteria Terpaginasi
// @Summary Ambil Data Kriteria (Pagination)
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Kata kunci pencarian"
// @Param pageSize query int false "Jumlah data per halaman"
// @Param pageNumber query int false "Nomor halaman yang diambil"
// @Param sortBy query string false "Parameter pengurutan"
// @Param sortType query string false "Tipe pengurutan [asc | desc]"
// @Success 200 {object} response.Base
// @Router /v1/kriteria/page [get]
func (h *KriteriaHandler) ResolvePaging(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	sortType := r.URL.Query().Get("sortType")

	if sortBy == "" {
		sortBy = "kode"
	}
	if sortType == "" {
		sortType = "asc"
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 10
	}

	pageNumber, _ := strconv.Atoi(pageNumberStr)
	if pageNumber <= 0 {
		pageNumber = 1
	}

	req := model.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
	}

	data, err := h.service.ResolvePaging(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// Update memperbarui data Kriteria berdasarkan ID.
// @Summary Perbarui data Kriteria
// @Tags Kriteria
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Param body body master.RequestKriteriaFormat true "Data Kriteria yang diperbarui"
// @Success 200 {object} response.Base
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
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Success 200 {object} response.Base
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
// @Tags Kriteria
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Kriteria"
// @Success 200 {object} response.Base
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