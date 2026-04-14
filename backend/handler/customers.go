package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"manufacture-backend/database"

	"github.com/golang-jwt/jwt/v5"
)

type customerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Password string `json:"password,omitempty"`
	Status   string `json:"status"`
}

type customerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type customerLoginResponse struct {
	Token string            `json:"token"`
	Data  database.Customer `json:"data"`
}

type customerResponse struct {
	Data database.Customer `json:"data"`
}

type customerListResponse struct {
	Data []database.Customer `json:"data"`
}

type customerOrdersResponse struct {
	Data []database.Order `json:"data"`
}

type customerClaims struct {
	CustomerID int64 `json:"customerId"`
	jwt.RegisteredClaims
}

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (h *Handler) Customers(w http.ResponseWriter, r *http.Request) {
	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/customers")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.handleListCustomers(w, r)
		case http.MethodPost:
			h.handleCreateCustomer(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetCustomer(w, r, id)
	case http.MethodPut:
		h.handleUpdateCustomer(w, r, id)
	case http.MethodDelete:
		h.handleDeleteCustomer(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListCustomers(w http.ResponseWriter, r *http.Request) {
	filter := database.CustomerFilter{}
	if search := strings.TrimSpace(r.URL.Query().Get("search")); search != "" {
		filter.Search = search
	}

	customers, err := database.ListCustomers(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil daftar pelanggan"})
		return
	}
	writeJSON(w, http.StatusOK, customerListResponse{Data: customers})
}

func (h *Handler) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email dan kata sandi pelanggan wajib diisi"})
		return
	}
	if !isValidEmail(req.Email) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email tidak valid"})
		return
	}
	if len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Kata sandi pelanggan minimal 8 karakter"})
		return
	}

	customer := database.Customer{
		Name:    strings.TrimSpace(req.Name),
		Email:   req.Email,
		Phone:   strings.TrimSpace(req.Phone),
		Address: strings.TrimSpace(req.Address),
		Status:  strings.TrimSpace(req.Status),
	}
	if req.Password != "" {
		hash, err := database.HashPassword(req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memproses kata sandi"})
			return
		}
		customer.PasswordHash = hash
	}

	id, err := database.CreateCustomer(h.db, customer)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat pelanggan"})
		return
	}

	created, err := database.GetCustomerByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pelanggan"})
		return
	}
	writeJSON(w, http.StatusCreated, customerResponse{Data: created})
}

func (h *Handler) handleGetCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	customer, err := database.GetCustomerByID(h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pelanggan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil pelanggan"})
		return
	}
	writeJSON(w, http.StatusOK, customerResponse{Data: customer})
}

func (h *Handler) handleUpdateCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email pelanggan wajib diisi"})
		return
	}
	if !isValidEmail(req.Email) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email tidak valid"})
		return
	}

	customer := database.Customer{
		ID:      id,
		Name:    strings.TrimSpace(req.Name),
		Email:   req.Email,
		Phone:   strings.TrimSpace(req.Phone),
		Address: strings.TrimSpace(req.Address),
		Status:  strings.TrimSpace(req.Status),
	}

	if req.Password != "" {
		if len(req.Password) < 8 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Kata sandi pelanggan minimal 8 karakter"})
			return
		}
		hash, err := database.HashPassword(req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memproses kata sandi"})
			return
		}
		customer.PasswordHash = hash
	}
	if err := database.UpdateCustomer(h.db, customer); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pelanggan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui pelanggan"})
		return
	}

	updated, err := database.GetCustomerByID(h.db, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat pelanggan"})
		return
	}
	writeJSON(w, http.StatusOK, customerResponse{Data: updated})
}

func (h *Handler) handleDeleteCustomer(w http.ResponseWriter, r *http.Request, id int64) {
	if err := database.DeleteCustomer(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pelanggan tidak ditemukan"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus pelanggan"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) CustomerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req customerLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email dan kata sandi wajib diisi"})
		return
	}

	customer, err := database.GetCustomerByEmail(h.db, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Email atau kata sandi salah"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memeriksa pelanggan"})
		return
	}

	if err := database.ComparePassword(customer.PasswordHash, req.Password); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Email atau kata sandi salah"})
		return
	}

	expiresIn, err := time.ParseDuration(h.cfg.CustomerTokenExpiry)
	if err != nil || expiresIn <= 0 {
		expiresIn = 24 * time.Hour
	}

	claims := customerClaims{
		CustomerID: customer.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "customer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.CustomerJWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat token"})
		return
	}

	writeJSON(w, http.StatusOK, customerLoginResponse{Token: signed, Data: customer})
}

func (h *Handler) CustomerOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerID, ok := h.verifyCustomerToken(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	filter := database.OrderFilter{CustomerID: customerID}
	orders, err := database.ListOrders(h.db, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil pesanan"})
		return
	}
	writeJSON(w, http.StatusOK, customerOrdersResponse{Data: orders})
}

func (h *Handler) CustomerProfile(w http.ResponseWriter, r *http.Request) {
	customerID, ok := h.verifyCustomerToken(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		customer, err := database.GetCustomerByID(h.db, customerID)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pelanggan tidak ditemukan"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil profil pelanggan"})
			return
		}
		writeJSON(w, http.StatusOK, customerResponse{Data: customer})
	case http.MethodPut:
		var req customerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)
		req.Phone = strings.TrimSpace(req.Phone)
		req.Address = strings.TrimSpace(req.Address)
		req.Password = strings.TrimSpace(req.Password)

		if req.Name == "" || req.Email == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Nama dan email wajib diisi"})
			return
		}
		if !isValidEmail(req.Email) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Email tidak valid"})
			return
		}
		if req.Password != "" && len(req.Password) < 8 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Kata sandi pelanggan minimal 8 karakter"})
			return
		}

		customer := database.Customer{
			ID:      customerID,
			Name:    req.Name,
			Email:   req.Email,
			Phone:   req.Phone,
			Address: req.Address,
			Status:  "active",
		}

		if req.Password != "" {
			hash, err := database.HashPassword(req.Password)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memproses kata sandi"})
				return
			}
			customer.PasswordHash = hash
		}

		if err := database.UpdateCustomer(h.db, customer); err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, map[string]string{"message": "Pelanggan tidak ditemukan"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui profil"})
			return
		}

		updated, err := database.GetCustomerByID(h.db, customerID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat profil pelanggan"})
			return
		}
		writeJSON(w, http.StatusOK, customerResponse{Data: updated})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) verifyCustomerToken(r *http.Request) (int64, bool) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return 0, false
	}

	const prefix = "Bearer "
	if len(authorization) <= len(prefix) || authorization[:len(prefix)] != prefix {
		return 0, false
	}

	tokenString := authorization[len(prefix):]
	claims := &customerClaims{}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.CustomerJWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return 0, false
	}

	return claims.CustomerID, true
}
