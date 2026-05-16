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
	KriteriaHandler        handlers.KriteriaHandler
	AuthHandler            handlers.AuthHandler
	RoleHandler            handlers.RoleHandler
	SumberDanaHandler      handlers.SumberDanaHandler
	PaguAnggaranHandler    handlers.PaguAnggaranHandler
	UsulanProyekHandler    handlers.UsulanProyekHandler
	PenilaianUsulanHandler handlers.PenilaianUsulanHandler
	PerankinganHandler     handlers.PerankinganHandler
	MenuHandler            handlers.MenuHandler // Pastikan ini terdaftar
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

		// ROUTE PUBLIC / ROUTE YANG MENGATUR MIDDLEWARE SENDIRI
		r.DomainHandlers.AuthHandler.Router(rc)
		
		// Daftarkan MenuHandler di sini, karena di dalam internal menu.go 
		// sudah ada rc.Group(func(protected chi.Router) { protected.Use(middleware.JWTProtected) ... })
		r.DomainHandlers.MenuHandler.Router(rc)

		// ROUTE PRIVATE (Dijaga oleh Satpam JWT)
		rc.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTProtected)

			r.DomainHandlers.KriteriaHandler.Router(protected)
			r.DomainHandlers.AuthHandler.UserRouter(protected)
			r.DomainHandlers.RoleHandler.Router(protected)
			r.DomainHandlers.SumberDanaHandler.Router(protected)
			r.DomainHandlers.PaguAnggaranHandler.Router(protected)
			r.DomainHandlers.UsulanProyekHandler.Router(protected)
			r.DomainHandlers.PenilaianUsulanHandler.Router(protected)
			r.DomainHandlers.PerankinganHandler.Router(protected)
		})
	})
}