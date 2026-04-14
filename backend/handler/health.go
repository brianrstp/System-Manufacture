package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"manufacture-backend/config"
)

type Handler struct {
	db  *sql.DB
	cfg config.Config
}

func New(db *sql.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/admin/login", h.AdminLogin)
	mux.HandleFunc("/api/admin/overview", h.AdminOverview)
	mux.HandleFunc("/api/products", h.Products)
	mux.HandleFunc("/api/products/", h.Products)
	mux.HandleFunc("/api/categories", h.Categories)
	mux.HandleFunc("/api/categories/", h.Categories)
	mux.HandleFunc("/api/units", h.Units)
	mux.HandleFunc("/api/units/", h.Units)
	mux.HandleFunc("/api/warehouses", h.Warehouses)
	mux.HandleFunc("/api/warehouses/", h.Warehouses)
	mux.HandleFunc("/api/boms", h.BOMs)
	mux.HandleFunc("/api/boms/", h.BOMs)
	mux.HandleFunc("/api/inventory", h.Inventories)
	mux.HandleFunc("/api/inventory/", h.Inventories)
	mux.HandleFunc("/api/stock_movements", h.StockMovements)
	mux.HandleFunc("/api/stock_movements/", h.StockMovements)
	mux.HandleFunc("/api/orders", h.Orders)
	mux.HandleFunc("/api/orders/", h.Orders)
	mux.HandleFunc("/api/customers", h.Customers)
	mux.HandleFunc("/api/customers/", h.Customers)
	mux.HandleFunc("/api/customers/login", h.CustomerLogin)
	mux.HandleFunc("/api/customer/orders", h.CustomerOrders)
	mux.HandleFunc("/api/customer/profile", h.CustomerProfile)
	mux.HandleFunc("/api/production", h.Production)
	mux.HandleFunc("/api/production/", h.Production)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var status int
	if err := h.db.QueryRow("SELECT 1").Scan(&status); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "database unavailable"})
		return
	}

	response := map[string]any{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
		"db":     status,
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
