package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"manufacture-backend/database"
)

type productRequest struct {
	SKU             string  `json:"sku"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	CategoryID      *int64  `json:"categoryId,omitempty"`
	UnitID          *int64  `json:"unitId,omitempty"`
	ProductType     string  `json:"productType"`
	StandardPrice   float64 `json:"standardPrice"`
	CostPrice       float64 `json:"costPrice"`
	LeadTimeDays    int     `json:"leadTimeDays"`
	MinOrderQty     int     `json:"minOrderQty"`
	ReorderPoint    int     `json:"reorderPoint"`
	LifecycleStatus string  `json:"lifecycleStatus"`
}

type productListResponse struct {
	Data []database.Product `json:"data"`
}

type productResponse struct {
	Data database.Product `json:"data"`
}

func (h *Handler) Products(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/products")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListProducts(w, r)
		case http.MethodPost:
			h.handleCreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetProduct(w, r, id)
	case http.MethodPut:
		h.handleUpdateProduct(w, r, id)
	case http.MethodDelete:
		h.handleDeleteProduct(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	filter := database.ProductFilter{}
	if search := strings.TrimSpace(r.URL.Query().Get("search")); search != "" {
		filter.Search = search
	}
	if categoryID := strings.TrimSpace(r.URL.Query().Get("categoryId")); categoryID != "" {
		id, err := strconv.ParseInt(categoryID, 10, 64)
		if err == nil {
			filter.CategoryID = &id
		}
	}
	filter.ProductType = strings.TrimSpace(r.URL.Query().Get("productType"))
	filter.LifecycleStatus = strings.TrimSpace(r.URL.Query().Get("lifecycleStatus"))

	products, err := database.ListProducts(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar produk"})
		return
	}

	writeJSON(w, http.StatusOK, productListResponse{Data: products})
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	product := database.Product{
		SKU:             strings.TrimSpace(req.SKU),
		Name:            strings.TrimSpace(req.Name),
		Description:     strings.TrimSpace(req.Description),
		CategoryID:      req.CategoryID,
		UnitID:          req.UnitID,
		ProductType:     strings.TrimSpace(req.ProductType),
		StandardPrice:   req.StandardPrice,
		CostPrice:       req.CostPrice,
		LeadTimeDays:    req.LeadTimeDays,
		MinOrderQty:     req.MinOrderQty,
		ReorderPoint:    req.ReorderPoint,
		LifecycleStatus: strings.TrimSpace(req.LifecycleStatus),
	}

	id, err := database.CreateProduct(h.db, product)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat produk"})
		return
	}

	created, err := database.GetProductByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat produk yang dibuat"})
		return
	}

	writeJSON(w, http.StatusCreated, productResponse{Data: created})
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request, id int64) {
	product, err := database.GetProductByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Produk tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil produk"})
		return
	}

	writeJSON(w, http.StatusOK, productResponse{Data: product})
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request, id int64) {
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	product := database.Product{
		ID:              id,
		SKU:             strings.TrimSpace(req.SKU),
		Name:            strings.TrimSpace(req.Name),
		Description:     strings.TrimSpace(req.Description),
		CategoryID:      req.CategoryID,
		UnitID:          req.UnitID,
		ProductType:     strings.TrimSpace(req.ProductType),
		StandardPrice:   req.StandardPrice,
		CostPrice:       req.CostPrice,
		LeadTimeDays:    req.LeadTimeDays,
		MinOrderQty:     req.MinOrderQty,
		ReorderPoint:    req.ReorderPoint,
		LifecycleStatus: strings.TrimSpace(req.LifecycleStatus),
	}

	if err := database.UpdateProduct(h.db, product); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Produk tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui produk"})
		return
	}

	updated, err := database.GetProductByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat produk yang diperbarui"})
		return
	}
	writeJSON(w, http.StatusOK, productResponse{Data: updated})
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteProduct(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Produk tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus produk"})
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
