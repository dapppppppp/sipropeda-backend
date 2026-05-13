package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"sipropeda-backend/configs"
	"sipropeda-backend/transport/http/router"
)

// HTTP struct menampung konfigurasi dan router
type HTTP struct {
	Config *configs.Config
	Router router.Router
}

// ProvideHTTP pembuat instance HTTP untuk Wire nanti
func ProvideHTTP(config *configs.Config, r router.Router) *HTTP {
	return &HTTP{
		Config: config,
		Router: r,
	}
}

// SetupAndServe menyalakan server agar terus standby
func (h *HTTP) SetupAndServe() {
	mux := chi.NewRouter()

	// 1. Pasang jalur router yang sudah kita buat
	h.Router.SetupRoutes(mux)

	// 2. Ambil port dari file .env (default 8080)
	port := h.Config.Server.Port
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("🚀 Server API SIPROPEDA standby dan siap menerima request di http://localhost%s\n", addr)
	fmt.Println("👉 Coba buka di browser: http://localhost" + addr + "/v1/ping")

	// 3. Jalankan server (kode akan "berhenti/standby" di baris ini)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("❌ Gagal menjalankan server HTTP: %v", err)
	}
}