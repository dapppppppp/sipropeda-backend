package handlers

import (
	"net/http"
	"strconv"

	"sipropeda-backend/internal/domain/transaction"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
)

type DashboardHandler struct {
	service transaction.DashboardService
}

func ProvideDashboardHandler(service transaction.DashboardService) DashboardHandler {
	return DashboardHandler{service: service}
}

// UBAH BAGIAN INI MENJADI LEBIH SEDERHANA
func (h *DashboardHandler) Router(r chi.Router) {
	r.Route("/dashboard", func(rc chi.Router) {
		rc.Get("/statistik", h.GetStatistik)
	})
}

func (h *DashboardHandler) GetStatistik(w http.ResponseWriter, r *http.Request) {
	tahunStr := r.URL.Query().Get("tahun")
	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		tahun = 2025 
	}

	data, err := h.service.GetStatistik(tahun)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}