package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type bomRequest struct {
	ProductID          int64   `json:"productId"`
	ComponentProductID int64   `json:"componentProductId"`
	ComponentUnitID    int64   `json:"componentUnitId"`
	ComponentQty       float64 `json:"componentQty"`
	WastePercentage    float64 `json:"wastePercentage"`
	ParentBOMID        *int64  `json:"parentBomId,omitempty"`
}

type bomResponse struct {
	Data database.BOMLine `json:"data"`
}

type bomListResponse struct {
	Data []database.BOMLine `json:"data"`
}

func (h *Handler) BOMs(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/boms")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListBOMs(w, r)
		case http.MethodPost:
			h.handleCreateBOM(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid BOM ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetBOM(w, r, id)
	case http.MethodPut:
		h.handleUpdateBOM(w, r, id)
	case http.MethodDelete:
		h.handleDeleteBOM(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListBOMs(w http.ResponseWriter, r *http.Request) {
	filter := database.BOMFilter{}
	if productID := strings.TrimSpace(r.URL.Query().Get("productId")); productID != "" {
		id, err := strconv.ParseInt(productID, 10, 64)
		if err == nil {
			filter.ProductID = &id
		}
	}
	bomLines, err := database.ListBOMLines(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar BOM"})
		return
	}
	writeJSON(w, http.StatusOK, bomListResponse{Data: bomLines})
}

func (h *Handler) handleCreateBOM(w http.ResponseWriter, r *http.Request) {
	var req bomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	bom := database.BOMLine{
		ProductID:          req.ProductID,
		ComponentProductID: req.ComponentProductID,
		ComponentUnitID:    req.ComponentUnitID,
		ComponentQty:       req.ComponentQty,
		WastePercentage:    req.WastePercentage,
		ParentBOMID:        req.ParentBOMID,
	}
	id, err := database.CreateBOMLine(h.db, bom)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat BOM"})
		return
	}
	created, err := database.GetBOMLineByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat BOM"})
		return
	}
	writeJSON(w, http.StatusCreated, bomResponse{Data: created})
}

func (h *Handler) handleGetBOM(w http.ResponseWriter, r *http.Request, id int64) {
	bom, err := database.GetBOMLineByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "BOM tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil BOM"})
		return
	}
	writeJSON(w, http.StatusOK, bomResponse{Data: bom})
}

func (h *Handler) handleUpdateBOM(w http.ResponseWriter, r *http.Request, id int64) {
	var req bomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	bom := database.BOMLine{
		ID:                 id,
		ProductID:          req.ProductID,
		ComponentProductID: req.ComponentProductID,
		ComponentUnitID:    req.ComponentUnitID,
		ComponentQty:       req.ComponentQty,
		WastePercentage:    req.WastePercentage,
		ParentBOMID:        req.ParentBOMID,
	}
	if err := database.UpdateBOMLine(h.db, bom); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "BOM tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui BOM"})
		return
	}
	updated, err := database.GetBOMLineByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat BOM"})
		return
	}
	writeJSON(w, http.StatusOK, bomResponse{Data: updated})
}

func (h *Handler) handleDeleteBOM(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteBOMLine(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "BOM tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus BOM"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
