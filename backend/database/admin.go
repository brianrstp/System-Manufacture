package database

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func EnsureAdminUser(db *sql.DB, username, password string) error {
	if err := createAdminTable(db); err != nil {
		return err
	}

	_, err := GetAdminByUsername(db, username)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO admins (username, password_hash, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		username,
		hash,
	)
	return err
}

func createAdminTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS admins (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}

func GetAdminByUsername(db *sql.DB, username string) (Admin, error) {
	var admin Admin
	row := db.QueryRow("SELECT id, username, password_hash, created_at, updated_at FROM admins WHERE username = ?", username)
	if err := row.Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.CreatedAt, &admin.UpdatedAt); err != nil {
		return Admin{}, err
	}
	return admin, nil
}

func HashPassword(password string) (string, error) {
	return hashPassword(password)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
