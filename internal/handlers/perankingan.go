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
)

type PerankinganHandler struct {
	service transaction.PerankinganService
}

func ProvidePerankinganHandler(service transaction.PerankinganService) PerankinganHandler {
	return PerankinganHandler{service: service}
}

func (h *PerankinganHandler) Router(r chi.Router) {
	r.Route("/perankingan", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Post("/hitung", h.HitungTOPSIS)
			protected.Get("/arsip", h.GetArsip)
		})
	})
}

func (h *PerankinganHandler) HitungTOPSIS(w http.ResponseWriter, r *http.Request) {
	var req transaction.RequestHitungTopsis
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	data, err := h.service.HitungTOPSIS(req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Perankingan TOPSIS berhasil dikalkulasi",
		"data":    data,
	})
}

// Menangkap param filter dan pagination
func (h *PerankinganHandler) GetArsip(w http.ResponseWriter, r *http.Request) {
	tahunStr := r.URL.Query().Get("tahun")
	tahap := r.URL.Query().Get("tahap")

	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		response.WithError(w, errors.New("parameter tahun tidak valid"))
		return
	}

	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	sortType := r.URL.Query().Get("sortType")

	// Default Value jika kosong
	if sortBy == "" {
		sortBy = "ranking"
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

	data, err := h.service.GetArsip(req, tahun, tahap)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}