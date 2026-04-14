package main

import (
	"log"
	"net/http"
	"time"

	"manufacture-backend/config"
	"manufacture-backend/database"
	"manufacture-backend/handler"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("gagal membuka database: %v", err)
	}
	defer db.Close()

	if err := database.EnsureAdminUser(db, cfg.AdminUser, cfg.AdminPassword); err != nil {
		log.Fatalf("gagal memastikan admin user: %v", err)
	}

	if err := database.EnsureManufacturingSchema(db); err != nil {
		log.Fatalf("gagal memastikan schema manufaktur: %v", err)
	}

	mux := http.NewServeMux()
	h := handler.New(db, cfg)
	handler.RegisterRoutes(mux, h)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Backend berjalan pada http://localhost:%s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
