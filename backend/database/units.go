package database

import (
	"database/sql"
	"time"
)

type Unit struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Factor      float64    `json:"factor"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

func CreateUnit(db *sql.DB, u Unit) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO units (code, name, description, factor, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		u.Code,
		u.Name,
		u.Description,
		u.Factor,
		u.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetUnitByID(db *sql.DB, id int64) (Unit, error) {
	var u Unit
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, code, name, description, factor, status, created_at, updated_at, deleted_at
		FROM units
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&u.ID,
		&u.Code,
		&u.Name,
		&u.Description,
		&u.Factor,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Unit{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		u.DeletedAt = &t
	}
	return u, nil
}

func ListUnits(db *sql.DB) ([]Unit, error) {
	rows, err := db.Query(`
		SELECT id, code, name, description, factor, status, created_at, updated_at, deleted_at
		FROM units
		WHERE deleted_at IS NULL
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := []Unit{}
	for rows.Next() {
		var u Unit
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&u.ID,
			&u.Code,
			&u.Name,
			&u.Description,
			&u.Factor,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			u.DeletedAt = &t
		}
		units = append(units, u)
	}

	return units, rows.Err()
}

func UpdateUnit(db *sql.DB, u Unit) error {
	result, err := db.Exec(`
		UPDATE units SET code = ?, name = ?, description = ?, factor = ?, status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		u.Code,
		u.Name,
		u.Description,
		u.Factor,
		u.Status,
		u.ID,
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

func DeleteUnit(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE units SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
