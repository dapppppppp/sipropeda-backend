package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/transport/http/middleware"
	"sipropeda-backend/transport/http/response"

	"github.com/go-chi/chi"
)

type AppConfigHandler struct {
	service auth.AppConfigService
}

func ProvideAppConfigHandler(service auth.AppConfigService) AppConfigHandler {
	return AppConfigHandler{service: service}
}

func (h *AppConfigHandler) Router(r chi.Router) {
	r.Route("/app-config", func(rc chi.Router) {
		
		// 1. Rute Publik (Tanpa Token JWT)
		rc.Get("/public/{id}", h.ResolveDTOByID)
		
		// 2. Rute File Server (Tanpa Token JWT, agar <v-img> bisa membacanya)
		rc.Get("/files", h.ServeFile)

		// 3. Rute Privat (Wajib Token JWT)
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected) // Sesuaikan dengan nama middleware Anda
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/{id}", h.Update)
			protected.Post("/upload", h.UploadFile)
		})
	})
}

// Menampilkan gambar yang telah diupload
func (h *AppConfigHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "File path is missing", http.StatusBadRequest)
		return
	}
	http.ServeFile(w, r, path)
}

func (h *AppConfigHandler) ResolveDTOByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WithError(w, errors.New("ID tidak boleh kosong"))
		return
	}

	data, err := h.service.ResolveDTOByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

func (h *AppConfigHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WithError(w, errors.New("ID tidak boleh kosong"))
		return
	}

	data, err := h.service.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

func (h *AppConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WithError(w, errors.New("ID tidak boleh kosong"))
		return
	}

	var req auth.RequestAppConfigFormat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, errors.New("invalid JSON body"))
		return
	}

	if err := h.service.Update(id, req); err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Konfigurasi Sistem berhasil diperbarui"})
}

func (h *AppConfigHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	path, err := h.service.UploadFile(r)
	if err != nil {
		response.WithError(w, err)
		return
	}
	
	// Kita return path stringnya secara langsung agar bisa dibaca res.data di Frontend
	response.WithJSON(w, http.StatusOK, path)
}