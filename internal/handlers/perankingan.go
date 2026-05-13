package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"sipropeda-backend/internal/domain/transaction"
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

// HitungTOPSIS mengeksekusi algoritma dan menyimpan ranking
// @Summary Kalkulasi Ranking TOPSIS
// @Description Endpoint untuk menghitung perangkingan menggunakan metode TOPSIS berdasarkan tahun dan tahap
// @Tags Perankingan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body transaction.RequestHitungTopsis true "Parameter Perhitungan"
// @Success 200 {object} response.Base
// @Router /v1/perankingan/hitung [post]
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

// GetArsip mengambil hasil ranking yang sudah tersimpan
// @Summary Ambil Arsip Perankingan
// @Tags Perankingan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param tahun query int true "Tahun Anggaran"
// @Param tahap query string true "Tahapan (misal: RKP)"
// @Success 200 {object} response.Base
// @Router /v1/perankingan/arsip [get]
func (h *PerankinganHandler) GetArsip(w http.ResponseWriter, r *http.Request) {
	tahunStr := r.URL.Query().Get("tahun")
	tahap := r.URL.Query().Get("tahap")

	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		response.WithError(w, errors.New("parameter tahun tidak valid"))
		return
	}

	data, err := h.service.GetArsip(tahun, tahap)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}