package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sipropeda-backend/internal/domain/master" // Sesuaikan dengan lokasi package model Anda
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
			protected.Post("/", h.Create)
			protected.Get("/{id}", h.ResolveByID)
			protected.Put("/{id}", h.Update)
			protected.Delete("/{id}", h.DeleteSoft)
		})
	})
}

// Create menambah data Bidang Pembangunan
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
func (h *BidangPembangunanHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ResolveAll()
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID mengambil data Bidang Pembangunan by ID
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