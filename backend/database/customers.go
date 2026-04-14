package database

import (
	"database/sql"
	"strings"
	"time"
)

type Customer struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Address      string     `json:"address"`
	PasswordHash string     `json:"-"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type CustomerFilter struct {
	Search string
}

func CreateCustomer(db *sql.DB, c Customer) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO customers (
			name, email, phone, address, password_hash, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		c.Name,
		c.Email,
		c.Phone,
		c.Address,
		c.PasswordHash,
		c.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetCustomerByID(db *sql.DB, id int64) (Customer, error) {
	var c Customer
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, name, email, phone, address, password_hash, status,
			created_at, updated_at, deleted_at
		FROM customers
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Address,
		&c.PasswordHash,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Customer{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		c.DeletedAt = &t
	}
	return c, nil
}

func GetCustomerByEmail(db *sql.DB, email string) (Customer, error) {
	var c Customer
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, name, email, phone, address, password_hash, status,
			created_at, updated_at, deleted_at
		FROM customers
		WHERE email = ? AND deleted_at IS NULL`, email)

	if err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Address,
		&c.PasswordHash,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Customer{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		c.DeletedAt = &t
	}
	return c, nil
}

func ListCustomers(db *sql.DB, filter CustomerFilter) ([]Customer, error) {
	query := `
		SELECT id, name, email, phone, address, status,
			created_at, updated_at, deleted_at
		FROM customers
		WHERE deleted_at IS NULL`
	args := []any{}

	if filter.Search != "" {
		query += ` AND (
			LOWER(name) LIKE ? OR
			LOWER(email) LIKE ? OR
			LOWER(phone) LIKE ? OR
			LOWER(address) LIKE ?
		)`
		pattern := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}
	query += " ORDER BY name ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var c Customer
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.Address,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			c.DeletedAt = &t
		}
		customers = append(customers, c)
	}

	return customers, rows.Err()
}

func CountCustomers(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func UpdateCustomer(db *sql.DB, c Customer) error {
	query := `UPDATE customers SET name = ?, email = ?, phone = ?, address = ?, status = ?, updated_at = NOW(3)`
	args := []any{c.Name, c.Email, c.Phone, c.Address, c.Status}
	if c.PasswordHash != "" {
		query += ", password_hash = ?"
		args = append(args, c.PasswordHash)
	}
	query += " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, c.ID)

	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteCustomer(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE customers SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
