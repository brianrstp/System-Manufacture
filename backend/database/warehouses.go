package database

import (
	"database/sql"
	"strings"
	"time"
)

type Warehouse struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type WarehouseFilter struct {
	Search string
	Status string
}

func CreateWarehouse(db *sql.DB, w Warehouse) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO warehouses (code, name, description, location, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		w.Code,
		w.Name,
		w.Description,
		w.Location,
		w.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetWarehouseByID(db *sql.DB, id int64) (Warehouse, error) {
	var w Warehouse
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, code, name, description, location, status, created_at, updated_at, deleted_at
		FROM warehouses
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&w.ID,
		&w.Code,
		&w.Name,
		&w.Description,
		&w.Location,
		&w.Status,
		&w.CreatedAt,
		&w.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Warehouse{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		w.DeletedAt = &t
	}
	return w, nil
}

func ListWarehouses(db *sql.DB, filter WarehouseFilter) ([]Warehouse, error) {
	query := `
		SELECT id, code, name, description, location, status, created_at, updated_at, deleted_at
		FROM warehouses
		WHERE deleted_at IS NULL`
	args := []any{}

	if filter.Search != "" {
		query += ` AND (
			LOWER(code) LIKE ? OR
			LOWER(name) LIKE ? OR
			LOWER(location) LIKE ? OR
			LOWER(status) LIKE ?
		)`
		pattern := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	query += " ORDER BY name ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	warehouses := []Warehouse{}
	for rows.Next() {
		var w Warehouse
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&w.ID,
			&w.Code,
			&w.Name,
			&w.Description,
			&w.Location,
			&w.Status,
			&w.CreatedAt,
			&w.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			w.DeletedAt = &t
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, rows.Err()
}

func UpdateWarehouse(db *sql.DB, w Warehouse) error {
	result, err := db.Exec(`
		UPDATE warehouses SET code = ?, name = ?, description = ?, location = ?, status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		w.Code,
		w.Name,
		w.Description,
		w.Location,
		w.Status,
		w.ID,
	)
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

func DeleteWarehouse(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE warehouses SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
