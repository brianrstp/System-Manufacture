package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"manufacture-backend/database"

	"github.com/golang-jwt/jwt/v5"
)

type adminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type adminLoginResponse struct {
	Token string `json:"token"`
}

type adminOverviewResponse struct {
	TotalProduction        int            `json:"totalProduction"`
	ActiveOrders           int            `json:"activeOrders"`
	TotalCustomers         int            `json:"totalCustomers"`
	TotalRevenue           int            `json:"totalRevenue"`
	MonthlyRevenue         []float64      `json:"monthlyRevenue"`
	ProductionStatusCounts map[string]int `json:"productionStatusCounts"`
	SupportTickets         int            `json:"supportTickets"`
	PendingApprovals       int            `json:"pendingApprovals"`
}

type adminClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req adminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	admin, err := database.GetAdminByUsername(h.db, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memeriksa akun"})
		return
	}

	if err := database.ComparePassword(admin.PasswordHash, req.Password); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		return
	}

	expiresIn, err := time.ParseDuration(h.cfg.AdminTokenExpiry)
	if err != nil || expiresIn <= 0 {
		expiresIn = time.Hour
	}

	claims := adminClaims{
		Username: admin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.AdminJWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal membuat token"})
		return
	}

	writeJSON(w, http.StatusOK, adminLoginResponse{Token: signed})
}

func (h *Handler) AdminOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.verifyAdminToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	activeOrders, err := database.CountOrders(h.db, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung pesanan aktif"})
		return
	}

	totalCustomers, err := database.CountCustomers(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung pelanggan"})
		return
	}

	totalRevenue, err := database.SumOrderAmounts(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung pendapatan"})
		return
	}

	totalProduction, err := database.CountProductionJobs(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung produksi"})
		return
	}

	monthlyRevenue, err := database.SumOrderAmountsByMonth(h.db, 6)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung pendapatan bulanan"})
		return
	}

	statusCounts, err := database.CountProductionJobsByStatus(h.db)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghitung status produksi"})
		return
	}

	writeJSON(w, http.StatusOK, adminOverviewResponse{
		TotalProduction:        totalProduction,
		ActiveOrders:           activeOrders,
		TotalCustomers:         totalCustomers,
		TotalRevenue:           int(totalRevenue),
		MonthlyRevenue:         monthlyRevenue,
		ProductionStatusCounts: statusCounts,
		SupportTickets:         7,
		PendingApprovals:       5,
	})
}

func (h *Handler) verifyAdminToken(r *http.Request) bool {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return false
	}

	const prefix = "Bearer "
	if len(authorization) <= len(prefix) || authorization[:len(prefix)] != prefix {
		return false
	}

	tokenString := authorization[len(prefix):]
	claims := &adminClaims{}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.AdminJWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return false
	}

	return claims.Username == h.cfg.AdminUser
}
