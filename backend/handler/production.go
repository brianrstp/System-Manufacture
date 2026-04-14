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

type productionJobRequest struct {
	JobCode      string `json:"jobCode"`
	Product      string `json:"product"`
	StartDate    string `json:"startDate"`
	DurationDays int    `json:"durationDays"`
	Status       string `json:"status"`
}

type productionJobResponse struct {
	Data database.ProductionJob `json:"data"`
}

type productionJobListResponse struct {
	Data []database.ProductionJob `json:"data"`
}

func (h *Handler) Production(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/production")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListProductionJobs(w, r)
		case http.MethodPost:
			h.handleCreateProductionJob(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid production job ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetProductionJob(w, r, id)
	case http.MethodPut:
		h.handleUpdateProductionJob(w, r, id)
	case http.MethodDelete:
		h.handleDeleteProductionJob(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListProductionJobs(w http.ResponseWriter, r *http.Request) {
	filter := database.ProductionJobFilter{}
	if search := strings.TrimSpace(r.URL.Query().Get("search")); search != "" {
		filter.Search = search
	}
	filter.Status = strings.TrimSpace(r.URL.Query().Get("status"))

	jobs, err := database.ListProductionJobs(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar produksi"})
		return
	}
	writeJSON(w, http.StatusOK, productionJobListResponse{Data: jobs})
}

func (h *Handler) handleCreateProductionJob(w http.ResponseWriter, r *http.Request) {
	var req productionJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now()
	}

	job := database.ProductionJob{
		JobCode:      strings.TrimSpace(req.JobCode),
		Product:      strings.TrimSpace(req.Product),
		StartDate:    startDate,
		DurationDays: req.DurationDays,
		Status:       strings.TrimSpace(req.Status),
	}

	id, err := database.CreateProductionJob(h.db, job)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat pekerjaan produksi"})
		return
	}

	created, err := database.GetProductionJobByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pekerjaan produksi"})
		return
	}
	writeJSON(w, http.StatusCreated, productionJobResponse{Data: created})
}

func (h *Handler) handleGetProductionJob(w http.ResponseWriter, r *http.Request, id int64) {
	job, err := database.GetProductionJobByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pekerjaan produksi tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil pekerjaan produksi"})
		return
	}
	writeJSON(w, http.StatusOK, productionJobResponse{Data: job})
}

func (h *Handler) handleUpdateProductionJob(w http.ResponseWriter, r *http.Request, id int64) {
	var req productionJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now()
	}

	job := database.ProductionJob{
		ID:           id,
		JobCode:      strings.TrimSpace(req.JobCode),
		Product:      strings.TrimSpace(req.Product),
		StartDate:    startDate,
		DurationDays: req.DurationDays,
		Status:       strings.TrimSpace(req.Status),
	}

	if err := database.UpdateProductionJob(h.db, job); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pekerjaan produksi tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui pekerjaan produksi"})
		return
	}

	updated, err := database.GetProductionJobByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pekerjaan produksi"})
		return
	}
	writeJSON(w, http.StatusOK, productionJobResponse{Data: updated})
}

func (h *Handler) handleDeleteProductionJob(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteProductionJob(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pekerjaan produksi tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus pekerjaan produksi"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
