package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"manufacture-backend/database"
)

type orderRequest struct {
	OrderNumber string  `json:"orderNumber"`
	CustomerID  int64   `json:"customerId"`
	Product     string  `json:"product"`
	OrderDate   string  `json:"orderDate"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
}

type orderResponse struct {
	Data database.Order `json:"data"`
}

type orderListResponse struct {
	Data []database.Order `json:"data"`
}

func (h *Handler) Orders(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/orders")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListOrders(w, r)
		case http.MethodPost:
			h.handleCreateOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetOrder(w, r, id)
	case http.MethodPut:
		h.handleUpdateOrder(w, r, id)
	case http.MethodDelete:
		h.handleDeleteOrder(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	filter := database.OrderFilter{}
	if search := strings.TrimSpace(r.URL.Query().Get("search")); search != "" {
		filter.Search = search
	}
	filter.Status = strings.TrimSpace(r.URL.Query().Get("status"))
	if limit := strings.TrimSpace(r.URL.Query().Get("limit")); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	orders, err := database.ListOrders(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar pesanan"})
		return
	}
	writeJSON(w, http.StatusOK, orderListResponse{Data: orders})
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		orderDate = time.Now()
	}

	order := database.Order{
		OrderNumber: strings.TrimSpace(req.OrderNumber),
		CustomerID:  req.CustomerID,
		Product:     strings.TrimSpace(req.Product),
		OrderDate:   orderDate,
		Amount:      req.Amount,
		Status:      strings.TrimSpace(req.Status),
	}

	id, err := database.CreateOrder(h.db, order)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat pesanan"})
		return
	}

	created, err := database.GetOrderByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pesanan"})
		return
	}
	writeJSON(w, http.StatusCreated, orderResponse{Data: created})
}

func (h *Handler) handleGetOrder(w http.ResponseWriter, r *http.Request, id int64) {
	order, err := database.GetOrderByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pesanan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil pesanan"})
		return
	}
	writeJSON(w, http.StatusOK, orderResponse{Data: order})
}

func (h *Handler) handleUpdateOrder(w http.ResponseWriter, r *http.Request, id int64) {
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		orderDate = time.Now()
	}

	order := database.Order{
		ID:          id,
		OrderNumber: strings.TrimSpace(req.OrderNumber),
		CustomerID:  req.CustomerID,
		Product:     strings.TrimSpace(req.Product),
		OrderDate:   orderDate,
		Amount:      req.Amount,
		Status:      strings.TrimSpace(req.Status),
	}

	if err := database.UpdateOrder(h.db, order); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pesanan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui pesanan"})
		return
	}

	updated, err := database.GetOrderByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pesanan"})
		return
	}
	writeJSON(w, http.StatusOK, orderResponse{Data: updated})
}

func (h *Handler) handleDeleteOrder(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteOrder(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pesanan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus pesanan"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
