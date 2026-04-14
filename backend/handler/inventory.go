package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type inventoryRequest struct {
	ProductID   int64   `json:"productId"`
	WarehouseID int64   `json:"warehouseId"`
	QtyOnHand   float64 `json:"qtyOnHand"`
	QtyReserved float64 `json:"qtyReserved"`
}

type inventoryResponse struct {
	Data database.Inventory `json:"data"`
}

type inventoryListResponse struct {
	Data []database.Inventory `json:"data"`
}

func (h *Handler) Inventories(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/inventory")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListInventory(w, r)
		case http.MethodPost:
			h.handleCreateInventory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetInventory(w, r, id)
	case http.MethodPut:
		h.handleUpdateInventory(w, r, id)
	case http.MethodDelete:
		h.handleDeleteInventory(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListInventory(w http.ResponseWriter, r *http.Request) {
	filter := database.InventoryFilter{}
	if productID := strings.TrimSpace(r.URL.Query().Get("productId")); productID != "" {
		id, err := strconv.ParseInt(productID, 10, 64)
		if err == nil {
			filter.ProductID = &id
		}
	}
	if warehouseID := strings.TrimSpace(r.URL.Query().Get("warehouseId")); warehouseID != "" {
		id, err := strconv.ParseInt(warehouseID, 10, 64)
		if err == nil {
			filter.WarehouseID = &id
		}
	}

	inventories, err := database.ListInventory(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil inventaris"})
		return
	}
	writeJSON(w, http.StatusOK, inventoryListResponse{Data: inventories})
}

func (h *Handler) handleCreateInventory(w http.ResponseWriter, r *http.Request) {
	var req inventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	inv := database.Inventory{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		QtyOnHand:   req.QtyOnHand,
		QtyReserved: req.QtyReserved,
	}
	id, err := database.CreateInventory(h.db, inv)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat inventaris"})
		return
	}
	created, err := database.GetInventoryByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat inventaris"})
		return
	}
	writeJSON(w, http.StatusCreated, inventoryResponse{Data: created})
}

func (h *Handler) handleGetInventory(w http.ResponseWriter, r *http.Request, id int64) {
	inv, err := database.GetInventoryByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Inventaris tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil inventaris"})
		return
	}
	writeJSON(w, http.StatusOK, inventoryResponse{Data: inv})
}

func (h *Handler) handleUpdateInventory(w http.ResponseWriter, r *http.Request, id int64) {
	var req inventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	inv := database.Inventory{
		ID:          id,
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		QtyOnHand:   req.QtyOnHand,
		QtyReserved: req.QtyReserved,
	}
	if err := database.UpdateInventory(h.db, inv); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Inventaris tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui inventaris"})
		return
	}
	updated, err := database.GetInventoryByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat inventaris"})
		return
	}
	writeJSON(w, http.StatusOK, inventoryResponse{Data: updated})
}

func (h *Handler) handleDeleteInventory(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteInventory(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Inventaris tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus inventaris"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
