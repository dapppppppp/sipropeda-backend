package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	_ "sipropeda-backend/docs"
	"sipropeda-backend/internal/handlers"
	"sipropeda-backend/transport/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

type DomainHandlers struct {
	KriteriaHandler          handlers.KriteriaHandler
	AuthHandler              handlers.AuthHandler
	RoleHandler              handlers.RoleHandler
	SumberDanaHandler        handlers.SumberDanaHandler
	BidangPembangunanHandler handlers.BidangPembangunanHandler
	PaguAnggaranHandler      handlers.PaguAnggaranHandler
	UsulanProyekHandler      handlers.UsulanProyekHandler
	PenilaianUsulanHandler   handlers.PenilaianUsulanHandler
	PerankinganHandler       handlers.PerankinganHandler
	MenuHandler              handlers.MenuHandler
	AppConfigHandler         handlers.AppConfigHandler
}

type Router struct {
	DomainHandlers DomainHandlers
}

func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

func (r *Router) SetupRoutes(mux *chi.Mux) {
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	mux.Get("/swagger/*", httpSwagger.WrapHandler)

	mux.Route("/v1", func(rc chi.Router) {
		rc.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "PONG!", "status": "success"}`))
		})

		// 👇 1. RUTE PUBLIK UNTUK MEMBACA GAMBAR (BEBAS DARI JWT) 👇
		rc.Get("/files", func(w http.ResponseWriter, req *http.Request) {
			filePath := req.URL.Query().Get("path")
			if filePath == "" {
				http.Error(w, "Path tidak ditemukan", http.StatusBadRequest)
				return
			}
			http.ServeFile(w, req, filePath)
		})

		// ROUTE PUBLIC / ROUTE YANG MENGATUR MIDDLEWARE SENDIRI
		r.DomainHandlers.AuthHandler.Router(rc)
		r.DomainHandlers.MenuHandler.Router(rc)
		
		// 👇 2. APP CONFIG DIPINDAHKAN KE SINI 👇
		r.DomainHandlers.AppConfigHandler.Router(rc)

		// ROUTE PRIVATE (Dijaga oleh Satpam JWT)
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)

			r.DomainHandlers.KriteriaHandler.Router(protected)
			r.DomainHandlers.AuthHandler.UserRouter(protected)
			r.DomainHandlers.RoleHandler.Router(protected)
			r.DomainHandlers.SumberDanaHandler.Router(protected)
			r.DomainHandlers.BidangPembangunanHandler.Router(protected)
			r.DomainHandlers.PaguAnggaranHandler.Router(protected)
			r.DomainHandlers.UsulanProyekHandler.Router(protected)
			r.DomainHandlers.PenilaianUsulanHandler.Router(protected)
			r.DomainHandlers.PerankinganHandler.Router(protected)
			
			// AppConfigHandler SUDAH DIHAPUS DARI SINI
		})
	})
}