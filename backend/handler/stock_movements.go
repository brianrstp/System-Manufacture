package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type stockMovementRequest struct {
	InventoryID   *int64  `json:"inventoryId,omitempty"`
	ProductID     int64   `json:"productId"`
	WarehouseID   int64   `json:"warehouseId"`
	MovementType  string  `json:"movementType"`
	Quantity      float64 `json:"quantity"`
	UnitID        *int64  `json:"unitId,omitempty"`
	ReferenceType string  `json:"referenceType"`
	ReferenceID   string  `json:"referenceId"`
	Notes         string  `json:"notes"`
}

type stockMovementResponse struct {
	Data database.StockMovement `json:"data"`
}

type stockMovementListResponse struct {
	Data []database.StockMovement `json:"data"`
}

func (h *Handler) StockMovements(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/stock_movements")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListStockMovements(w, r)
		case http.MethodPost:
			h.handleCreateStockMovement(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid stock movement ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetStockMovement(w, r, id)
	case http.MethodPut:
		h.handleUpdateStockMovement(w, r, id)
	case http.MethodDelete:
		h.handleDeleteStockMovement(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListStockMovements(w http.ResponseWriter, r *http.Request) {
	filter := database.StockMovementFilter{}
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
	filter.MovementType = strings.TrimSpace(r.URL.Query().Get("movementType"))

	movements, err := database.ListStockMovements(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil stock movement"})
		return
	}
	writeJSON(w, http.StatusOK, stockMovementListResponse{Data: movements})
}

func (h *Handler) handleCreateStockMovement(w http.ResponseWriter, r *http.Request) {
	var req stockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	movement := database.StockMovement{
		InventoryID:   req.InventoryID,
		ProductID:     req.ProductID,
		WarehouseID:   req.WarehouseID,
		MovementType:  strings.TrimSpace(req.MovementType),
		Quantity:      req.Quantity,
		UnitID:        req.UnitID,
		ReferenceType: strings.TrimSpace(req.ReferenceType),
		ReferenceID:   strings.TrimSpace(req.ReferenceID),
		Notes:         strings.TrimSpace(req.Notes),
	}
	id, err := database.CreateStockMovement(h.db, movement)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat stock movement"})
		return
	}
	created, err := database.GetStockMovementByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat stock movement"})
		return
	}
	writeJSON(w, http.StatusCreated, stockMovementResponse{Data: created})
}

func (h *Handler) handleGetStockMovement(w http.ResponseWriter, r *http.Request, id int64) {
	movement, err := database.GetStockMovementByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Stock movement tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil stock movement"})
		return
	}
	writeJSON(w, http.StatusOK, stockMovementResponse{Data: movement})
}

func (h *Handler) handleUpdateStockMovement(w http.ResponseWriter, r *http.Request, id int64) {
	var req stockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	movement := database.StockMovement{
		ID:            id,
		InventoryID:   req.InventoryID,
		ProductID:     req.ProductID,
		WarehouseID:   req.WarehouseID,
		MovementType:  strings.TrimSpace(req.MovementType),
		Quantity:      req.Quantity,
		UnitID:        req.UnitID,
		ReferenceType: strings.TrimSpace(req.ReferenceType),
		ReferenceID:   strings.TrimSpace(req.ReferenceID),
		Notes:         strings.TrimSpace(req.Notes),
	}
	if err := database.UpdateStockMovement(h.db, movement); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Stock movement tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui stock movement"})
		return
	}
	updated, err := database.GetStockMovementByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat stock movement"})
		return
	}
	writeJSON(w, http.StatusOK, stockMovementResponse{Data: updated})
}

func (h *Handler) handleDeleteStockMovement(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteStockMovement(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Stock movement tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus stock movement"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
