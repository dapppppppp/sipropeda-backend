package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/transaction"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
)

type PenilaianUsulanHandler struct {
	service transaction.PenilaianUsulanService
}

func ProvidePenilaianUsulanHandler(service transaction.PenilaianUsulanService) PenilaianUsulanHandler {
	return PenilaianUsulanHandler{service: service}
}

func (h *PenilaianUsulanHandler) Router(r chi.Router) {
	r.Route("/penilaian-usulan", func(rc chi.Router) {
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)
			protected.Get("/{usulanId}", h.GetByUsulanID)
			protected.Post("/", h.SaveBulk)
		})
	})
}

// SaveBulk menyimpan atau memperbarui daftar nilai untuk satu usulan
// @Summary Simpan Penilaian Usulan (Bulk)
// @Tags Penilaian Usulan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body transaction.RequestBulkPenilaian true "Data Penilaian"
// @Success 200 {object} response.Base
// @Router /v1/penilaian-usulan [post]
func (h *PenilaianUsulanHandler) SaveBulk(w http.ResponseWriter, r *http.Request) {
	var req transaction.RequestBulkPenilaian
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.service.SaveBulk(req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Penilaian successfully saved"})
}

// GetByUsulanID mengambil daftar nilai berdasarkan Usulan ID
// @Summary Ambil Nilai berdasarkan Usulan ID
// @Tags Penilaian Usulan
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param usulanId path string true "ID Usulan Proyek"
// @Success 200 {object} response.Base
// @Router /v1/penilaian-usulan/{usulanId} [get]
func (h *PenilaianUsulanHandler) GetByUsulanID(w http.ResponseWriter, r *http.Request) {
	usulanID := chi.URLParam(r, "usulanId")

	data, err := h.service.GetByUsulanID(usulanID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}