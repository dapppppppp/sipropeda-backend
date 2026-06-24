package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"sipropeda-backend/internal/domain/transaction"
	"sipropeda-backend/shared/model"
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
			protected.Get("/all", h.GetAllData) // Route baru untuk Get All (tanpa pagination)
			protected.Get("/tahun-anggaran", h.GetAvailableYears)
			protected.Post("/", h.Create)
			protected.Post("/import", h.ImportExcel)
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/bulk-update-status", h.BulkUpdateStatus)
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

// ResolveAll mengambil semua data Usulan Proyek dengan filter dan pagination
// @Summary Ambil semua data Usulan Proyek
// @Tags Usulan Proyek
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Kata kunci pencarian"
// @Param pageSize query int false "Jumlah data per halaman"
// @Param pageNumber query int false "Nomor halaman yang diambil"
// @Param sortBy query string false "Parameter pengurutan"
// @Param sortType query string false "Tipe pengurutan [asc | desc]"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek [get]
func (h *UsulanProyekHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	sortType := r.URL.Query().Get("sortType")
	bidangId := r.URL.Query().Get("bidangId")
	sumberDanaId := r.URL.Query().Get("sumberDanaId")

	// Set Default Values
	if sortBy == "" {
		sortBy = "createdAt"
	}
	if sortType == "" {
		sortType = "desc"
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 10
	}

	pageNumber, _ := strconv.Atoi(pageNumberStr)
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// Masukkan ke StandardRequest
	req := model.StandardRequest{
		Keyword:      keyword,
		PageSize:     pageSize,
		PageNumber:   pageNumber,
		SortBy:       sortBy,
		SortType:     sortType,
		BidangId:     bidangId,
		SumberDanaId: sumberDanaId,
		Tahun:        r.URL.Query().Get("tahun"),
	}

	// Panggil Service
	data, err := h.service.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// GetAllData mengambil seluruh data Usulan Proyek tanpa pagination
// @Summary Ambil semua data Usulan Proyek (tanpa pagination)
// @Tags Usulan Proyek
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/all [get]
func (h *UsulanProyekHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAllData()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// GetAvailableYears mengambil daftar tahun anggaran yang tersedia
// @Summary Mengambil daftar tahun anggaran yang tersedia
// @Tags Usulan Proyek
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/tahun-anggaran [get]
func (h *UsulanProyekHandler) GetAvailableYears(w http.ResponseWriter, r *http.Request) {
	years, err := h.service.GetAvailableYears()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, years)
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

// BulkUpdateStatus mengubah status beberapa Usulan Proyek sekaligus
// @Summary Update status masal Usulan Proyek
// @Tags Usulan Proyek
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body transaction.RequestBulkUpdateStatus true "Data yang akan diedit masal"
// @Success 200 {object} response.Base
// @Router /v1/usulan-proyek/bulk-update-status [put]
func (h *UsulanProyekHandler) BulkUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req transaction.RequestBulkUpdateStatus
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

	if err := h.service.BulkUpdateStatus(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Status Usulan Proyek successfully updated in bulk"})
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

// ImportExcel memproses unggahan file RKPDes
// @Summary Import data Usulan Proyek dari Excel
// @Tags Usulan Proyek
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param file_excel formData file true "File Excel"
// @Param tahun_anggaran formData int true "Tahun Anggaran"
// @Success 201 {object} response.Base
// @Router /v1/usulan-proyek/import [post]
func (h *UsulanProyekHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	// Batasi ukuran memori multipart hingga 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.WithError(w, errors.New("ukuran file terlalu besar"))
		return
	}

	tahunStr := r.FormValue("tahun_anggaran")
	tahun, err := strconv.Atoi(tahunStr)
	if err != nil || tahun == 0 {
		response.WithError(w, errors.New("tahun anggaran tidak valid"))
		return
	}

	file, _, err := r.FormFile("file_excel")
	if err != nil {
		response.WithError(w, errors.New("file excel wajib diunggah"))
		return
	}
	defer file.Close()

	jumlahData, err := h.service.ImportExcelRKP(file, tahun)
	if err != nil {
		response.WithError(w, err)
		return
	}

	pesan := map[string]interface{}{
		"message": "Berhasil mengimpor " + strconv.Itoa(jumlahData) + " usulan proyek dari file RKPDes",
		"jumlah":  jumlahData,
	}

	response.WithJSON(w, http.StatusCreated, pesan)
}