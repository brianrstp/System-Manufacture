package database

import (
	"database/sql"
	"time"
)

type Category struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	ParentID    *int64     `json:"parentId,omitempty"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

func CreateCategory(db *sql.DB, c Category) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO categories (name, slug, parent_id, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		c.Name,
		c.Slug,
		nullableInt64Value(c.ParentID),
		c.Description,
		c.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetCategoryByID(db *sql.DB, id int64) (Category, error) {
	var c Category
	var parentID sql.NullInt64
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, name, slug, parent_id, description, status, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Slug,
		&parentID,
		&c.Description,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Category{}, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.Int64
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		c.DeletedAt = &t
	}
	return c, nil
}

func ListCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query(`
		SELECT id, name, slug, parent_id, description, status, created_at, updated_at, deleted_at
		FROM categories
		WHERE deleted_at IS NULL
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		var parentID sql.NullInt64
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Slug,
			&parentID,
			&c.Description,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.Int64
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			c.DeletedAt = &t
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func UpdateCategory(db *sql.DB, c Category) error {
	result, err := db.Exec(`
		UPDATE categories SET name = ?, slug = ?, parent_id = ?, description = ?, status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		c.Name,
		c.Slug,
		nullableInt64Value(c.ParentID),
		c.Description,
		c.Status,
		c.ID,
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

func DeleteCategory(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE categories SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
