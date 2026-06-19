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

type BidangPembangunanHandler struct {
	service master.BidangPembangunanService
}

func ProvideBidangPembangunanHandler(service master.BidangPembangunanService) BidangPembangunanHandler {
	return BidangPembangunanHandler{service: service}
}

func (h *BidangPembangunanHandler) Router(r chi.Router) {
	r.Route("/bidang-pembangunan", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/", h.ResolveAll)
			protected.Get("/page", h.ResolvePaging) // Endpoint khusus pagination tabel
			protected.Post("/", h.Create)
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/{id}", h.Update)
			protected.Delete("/{id}", h.DeleteSoft)
		})
	})
}

// Create menambah data Bidang Pembangunan
// @Summary Tambah data Bidang Pembangunan
// @Tags Bidang Pembangunan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body master.RequestBidangPembangunan true "Data Bidang Pembangunan"
// @Success 201 {object} response.Base
// @Router /v1/bidang-pembangunan [post]
func (h *BidangPembangunanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestBidangPembangunan
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

// ResolveAll mengambil semua data Bidang Pembangunan
// @Summary Ambil semua data Bidang Pembangunan
// @Tags Bidang Pembangunan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} response.Base
// @Router /v1/bidang-pembangunan [get]
func (h *BidangPembangunanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolvePaging mengambil data Bidang Pembangunan secara terpaginasi
// @Summary Ambil Data Bidang Pembangunan (Pagination)
// @Tags Bidang Pembangunan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param q query string false "Kata kunci pencarian"
// @Param pageSize query int false "Jumlah data per halaman"
// @Param pageNumber query int false "Nomor halaman yang diambil"
// @Param sortBy query string false "Parameter pengurutan"
// @Param sortType query string false "Tipe pengurutan [asc | desc]"
// @Success 200 {object} response.Base
// @Router /v1/bidang-pembangunan/page [get]
func (h *BidangPembangunanHandler) ResolvePaging(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	sortType := r.URL.Query().Get("sortType")

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

// ResolveByID mengambil data Bidang Pembangunan by ID
// @Summary Ambil detail Bidang Pembangunan by ID
// @Tags Bidang Pembangunan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Bidang Pembangunan"
// @Success 200 {object} response.Base
// @Router /v1/bidang-pembangunan/{id} [get]
func (h *BidangPembangunanHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
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

// Update mengubah data Bidang Pembangunan
// @Summary Update data Bidang Pembangunan
// @Tags Bidang Pembangunan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Bidang Pembangunan"
// @Param body body master.RequestBidangPembangunan true "Data yang akan diedit"
// @Success 200 {object} response.Base
// @Router /v1/bidang-pembangunan/{id} [put]
func (h *BidangPembangunanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req master.RequestBidangPembangunan
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
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Bidang Pembangunan successfully updated"})
}

// DeleteSoft menghapus data Bidang Pembangunan
// @Summary Hapus data Bidang Pembangunan (Soft Delete)
// @Tags Bidang Pembangunan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "ID Bidang Pembangunan"
// @Success 200 {object} response.Base
// @Router /v1/bidang-pembangunan/{id} [delete]
func (h *BidangPembangunanHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
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
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Bidang Pembangunan successfully deleted"})
}