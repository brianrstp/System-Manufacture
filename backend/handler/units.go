package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type unitRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Factor      float64 `json:"factor"`
	Status      string  `json:"status"`
}

type unitResponse struct {
	Data database.Unit `json:"data"`
}

type unitListResponse struct {
	Data []database.Unit `json:"data"`
}

func (h *Handler) Units(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/units")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListUnits(w, r)
		case http.MethodPost:
			h.handleCreateUnit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid unit ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetUnit(w, r, id)
	case http.MethodPut:
		h.handleUpdateUnit(w, r, id)
	case http.MethodDelete:
		h.handleDeleteUnit(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListUnits(w http.ResponseWriter, r *http.Request) {
	units, err := database.ListUnits(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar unit"})
		return
	}
	writeJSON(w, http.StatusOK, unitListResponse{Data: units})
}

func (h *Handler) handleCreateUnit(w http.ResponseWriter, r *http.Request) {
	var req unitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	unit := database.Unit{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Factor:      req.Factor,
		Status:      strings.TrimSpace(req.Status),
	}
	id, err := database.CreateUnit(h.db, unit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat unit"})
		return
	}
	created, err := database.GetUnitByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat unit"})
		return
	}
	writeJSON(w, http.StatusCreated, unitResponse{Data: created})
}

func (h *Handler) handleGetUnit(w http.ResponseWriter, r *http.Request, id int64) {
	unit, err := database.GetUnitByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Unit tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil unit"})
		return
	}
	writeJSON(w, http.StatusOK, unitResponse{Data: unit})
}

func (h *Handler) handleUpdateUnit(w http.ResponseWriter, r *http.Request, id int64) {
	var req unitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	unit := database.Unit{
		ID:          id,
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Factor:      req.Factor,
		Status:      strings.TrimSpace(req.Status),
	}
	if err := database.UpdateUnit(h.db, unit); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Unit tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui unit"})
		return
	}
	updated, err := database.GetUnitByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat unit"})
		return
	}
	writeJSON(w, http.StatusOK, unitResponse{Data: updated})
}

func (h *Handler) handleDeleteUnit(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteUnit(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Unit tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus unit"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
