.PHONY: dev swagger run

# Menjalankan server dengan fitur Live Reload (otomatis restart kalau kode disave)
dev:
	go run github.com/air-verse/air@latest

# Mengupdate dokumentasi Swagger setiap kali ada perubahan endpoint
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init

# Menjalankan server secara manual (fallback)
run:
	go run main.go