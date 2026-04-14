package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type categoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *int64 `json:"parentId,omitempty"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type categoryResponse struct {
	Data database.Category `json:"data"`
}

type categoryListResponse struct {
	Data []database.Category `json:"data"`
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/categories")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListCategories(w, r)
		case http.MethodPost:
			h.handleCreateCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetCategory(w, r, id)
	case http.MethodPut:
		h.handleUpdateCategory(w, r, id)
	case http.MethodDelete:
		h.handleDeleteCategory(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := database.ListCategories(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar kategori"})
		return
	}
	writeJSON(w, http.StatusOK, categoryListResponse{Data: categories})
}

func (h *Handler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	cat := database.Category{
		Name:        strings.TrimSpace(req.Name),
		Slug:        strings.TrimSpace(req.Slug),
		ParentID:    req.ParentID,
		Description: strings.TrimSpace(req.Description),
		Status:      strings.TrimSpace(req.Status),
	}
	id, err := database.CreateCategory(h.db, cat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat kategori"})
		return
	}
	created, err := database.GetCategoryByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat kategori"})
		return
	}
	writeJSON(w, http.StatusCreated, categoryResponse{Data: created})
}

func (h *Handler) handleGetCategory(w http.ResponseWriter, r *http.Request, id int64) {
	cat, err := database.GetCategoryByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kategori tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil kategori"})
		return
	}
	writeJSON(w, http.StatusOK, categoryResponse{Data: cat})
}

func (h *Handler) handleUpdateCategory(w http.ResponseWriter, r *http.Request, id int64) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}
	cat := database.Category{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Slug:        strings.TrimSpace(req.Slug),
		ParentID:    req.ParentID,
		Description: strings.TrimSpace(req.Description),
		Status:      strings.TrimSpace(req.Status),
	}
	if err := database.UpdateCategory(h.db, cat); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kategori tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui kategori"})
		return
	}
	updated, err := database.GetCategoryByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat kategori"})
		return
	}
	writeJSON(w, http.StatusOK, categoryResponse{Data: updated})
}

func (h *Handler) handleDeleteCategory(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteCategory(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kategori tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus kategori"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
