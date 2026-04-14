package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type warehouseRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Status      string `json:"status"`
}

type warehouseResponse struct {
	Data database.Warehouse `json:"data"`
}

type warehouseListResponse struct {
	Data []database.Warehouse `json:"data"`
}

func (h *Handler) Warehouses(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/warehouses")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListWarehouses(w, r)
		case http.MethodPost:
			h.handleCreateWarehouse(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetWarehouse(w, r, id)
	case http.MethodPut:
		h.handleUpdateWarehouse(w, r, id)
	case http.MethodDelete:
		h.handleDeleteWarehouse(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListWarehouses(w http.ResponseWriter, r *http.Request) {
	filter := database.WarehouseFilter{
		Search: strings.TrimSpace(r.URL.Query().Get("search")),
		Status: strings.TrimSpace(r.URL.Query().Get("status")),
	}
	warehouses, err := database.ListWarehouses(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar gudang"})
		return
	}
	writeJSON(w, http.StatusOK, warehouseListResponse{Data: warehouses})
}

func (h *Handler) handleCreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req warehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	warehouse := database.Warehouse{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Location:    strings.TrimSpace(req.Location),
		Status:      strings.TrimSpace(req.Status),
	}
	id, err := database.CreateWarehouse(h.db, warehouse)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat gudang"})
		return
	}
	created, err := database.GetWarehouseByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat gudang"})
		return
	}
	writeJSON(w, http.StatusCreated, warehouseResponse{Data: created})
}

func (h *Handler) handleGetWarehouse(w http.ResponseWriter, r *http.Request, id int64) {
	warehouse, err := database.GetWarehouseByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Gudang tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil gudang"})
		return
	}
	writeJSON(w, http.StatusOK, warehouseResponse{Data: warehouse})
}

func (h *Handler) handleUpdateWarehouse(w http.ResponseWriter, r *http.Request, id int64) {
	var req warehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	warehouse := database.Warehouse{
		ID:          id,
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Location:    strings.TrimSpace(req.Location),
		Status:      strings.TrimSpace(req.Status),
	}
	if err := database.UpdateWarehouse(h.db, warehouse); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Gudang tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui gudang"})
		return
	}
	updated, err := database.GetWarehouseByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat gudang"})
		return
	}
	writeJSON(w, http.StatusOK, warehouseResponse{Data: updated})
}

func (h *Handler) handleDeleteWarehouse(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteWarehouse(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Gudang tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus gudang"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
