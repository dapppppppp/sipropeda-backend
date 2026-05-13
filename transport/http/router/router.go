package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"sipropeda-backend/internal/handlers"
	"sipropeda-backend/transport/http/middleware"
	_ "sipropeda-backend/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type DomainHandlers struct {
	KriteriaHandler handlers.KriteriaHandler
	AuthHandler     handlers.AuthHandler 
	RoleHandler     handlers.RoleHandler // <-- Tambahan Role Handler
	SumberDanaHandler handlers.SumberDanaHandler // <-- Tambahan ini
	PaguAnggaranHandler handlers.PaguAnggaranHandler // <-- Tambahan
	UsulanProyekHandler handlers.UsulanProyekHandler // <-- Tambahan untuk Usulan Proyek
	PenilaianUsulanHandler handlers.PenilaianUsulanHandler // <-- Tambahan untuk Penilaian Usulan
	PerankinganHandler handlers.PerankinganHandler // <-- Tambahan untuk Perankingan
	MenuHandler handlers.MenuHandler // <-- Tambahan untuk Menu
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

		// ROUTE PUBLIC
		r.DomainHandlers.AuthHandler.Router(rc)

		// ROUTE PRIVATE (Dijaga oleh Satpam JWT)
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected) 
			
			r.DomainHandlers.KriteriaHandler.Router(protected)
			r.DomainHandlers.AuthHandler.UserRouter(protected) 
			r.DomainHandlers.RoleHandler.Router(protected) // <-- Daftarkan Router Role
			r.DomainHandlers.SumberDanaHandler.Router(protected) // <-- Tambahan ini
			r.DomainHandlers.PaguAnggaranHandler.Router(protected)
			r.DomainHandlers.UsulanProyekHandler.Router(protected)
			r.DomainHandlers.PenilaianUsulanHandler.Router(protected)
			r.DomainHandlers.PerankinganHandler.Router(protected)
			r.DomainHandlers.MenuHandler.Router(protected)
		})
	})
}