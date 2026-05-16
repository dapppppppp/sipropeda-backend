package main

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sipropeda-backend/configs"
	"sipropeda-backend/infras"
	"sipropeda-backend/internal/domain/auth"
	"sipropeda-backend/internal/domain/master"
	"sipropeda-backend/internal/domain/transaction"
	"sipropeda-backend/internal/handlers"
	"sipropeda-backend/transport/http"
	"sipropeda-backend/transport/http/router"
)

// @title SIPROPEDA API
// @version 1.0
// @description API Backend untuk Sistem Pendukung Keputusan (SIPROPEDA) - Metode TOPSIS.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config := configs.Get()

	dbConn := infras.ProvidePostgreSQLConn(config)
	err := dbConn.Read.Ping()
	if err != nil {
		fmt.Println("❌ Gagal terhubung ke Database:", err)
		return
	}
	fmt.Println("✅ SUKSES! Terhubung ke Database SIPROPEDA!")

	// -- Domain Kriteria
	kriteriaRepo := master.ProvideKriteriaRepository(dbConn)
	kriteriaSvc := master.ProvideKriteriaService(kriteriaRepo)
	kriteriaHdl := handlers.ProvideKriteriaHandler(kriteriaSvc)
	
	// -- Domain Master (Sumber Dana)
	sumberDanaRepo := master.ProvideSumberDanaRepository(dbConn)
	sumberDanaSvc := master.ProvideSumberDanaService(sumberDanaRepo)
	sumberDanaHdl := handlers.ProvideSumberDanaHandler(sumberDanaSvc)
	
	// -- Domain Auth (Role)
	roleRepo := auth.ProvideRoleRepository(dbConn)
	roleSvc := auth.ProvideRoleService(roleRepo)
	roleHdl := handlers.ProvideRoleHandler(roleSvc)
	
	// -- Domain Transaction (Pagu Anggaran)
	paguAnggaranRepo := transaction.ProvidePaguAnggaranRepository(dbConn)
	paguAnggaranSvc := transaction.ProvidePaguAnggaranService(paguAnggaranRepo)
	paguAnggaranHdl := handlers.ProvidePaguAnggaranHandler(paguAnggaranSvc)

	usulanProyekRepo := transaction.ProvideUsulanProyekRepository(dbConn)
	usulanProyekSvc := transaction.ProvideUsulanProyekService(usulanProyekRepo)
	usulanProyekHdl := handlers.ProvideUsulanProyekHandler(usulanProyekSvc)

	// -- Domain Transaction (Perankingan TOPSIS)
	perankinganRepo := transaction.ProvidePerankinganRepository(dbConn)
	perankinganSvc := transaction.ProvidePerankinganService(perankinganRepo)
	perankinganHdl := handlers.ProvidePerankinganHandler(perankinganSvc)

	// -- Domain Transaction (Penilaian Usulan)
	penilaianRepo := transaction.ProvidePenilaianUsulanRepository(dbConn)
	penilaianSvc := transaction.ProvidePenilaianUsulanService(penilaianRepo)
	penilaianHdl := handlers.ProvidePenilaianUsulanHandler(penilaianSvc)

	// -- Domain Auth (User)
	userRepo := auth.ProvideUserRepository(dbConn)
	userSvc := auth.ProvideUserService(userRepo)
	authHdl := handlers.ProvideAuthHandler(userSvc)

	// -- Domain Auth (Menu) <-- Pastikan memanggil nama function yang baru
	menuRepo := auth.ProvideMenuRepositoryPostgreSQL(dbConn) // <-- PERUBAHAN DI SINI
	menuSvc := auth.ProvideMenuServiceImpl(menuRepo)         // <-- PERUBAHAN DI SINI
	menuHdl := handlers.ProvideMenuHandler(menuSvc)

	// Masukkan Handler ke Router
	domainHandlers := router.DomainHandlers{
		KriteriaHandler:        kriteriaHdl,
		AuthHandler:            authHdl,
		RoleHandler:            roleHdl,
		SumberDanaHandler:      sumberDanaHdl,
		PaguAnggaranHandler:    paguAnggaranHdl,
		UsulanProyekHandler:    usulanProyekHdl,
		PenilaianUsulanHandler: penilaianHdl,
		PerankinganHandler:     perankinganHdl,
		MenuHandler:            menuHdl,
	}

	appRouter := router.ProvideRouter(domainHandlers)
	httpServer := http.ProvideHTTP(config, appRouter)

	go httpServer.SetupAndServe()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("\n🛑 Mematikan Server SIPROPEDA...")
}