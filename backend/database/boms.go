package database

import (
	"database/sql"
	"time"
)

type BOMLine struct {
	ID                 int64      `json:"id"`
	ProductID          int64      `json:"productId"`
	ComponentProductID int64      `json:"componentProductId"`
	ComponentUnitID    int64      `json:"componentUnitId"`
	ComponentQty       float64    `json:"componentQty"`
	WastePercentage    float64    `json:"wastePercentage"`
	ParentBOMID        *int64     `json:"parentBomId,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}

type BOMFilter struct {
	ProductID *int64
}

func CreateBOMLine(db *sql.DB, b BOMLine) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO bills_of_materials (
			product_id, component_product_id, component_unit_id, component_qty,
			waste_percentage, parent_bom_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		b.ProductID,
		b.ComponentProductID,
		b.ComponentUnitID,
		b.ComponentQty,
		b.WastePercentage,
		nullableInt64Value(b.ParentBOMID),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetBOMLineByID(db *sql.DB, id int64) (BOMLine, error) {
	var b BOMLine
	var parentID sql.NullInt64
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, product_id, component_product_id, component_unit_id,
			component_qty, waste_percentage, parent_bom_id, created_at, updated_at, deleted_at
		FROM bills_of_materials
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&b.ID,
		&b.ProductID,
		&b.ComponentProductID,
		&b.ComponentUnitID,
		&b.ComponentQty,
		&b.WastePercentage,
		&parentID,
		&b.CreatedAt,
		&b.UpdatedAt,
		&deletedAt,
	); err != nil {
		return BOMLine{}, err
	}
	if parentID.Valid {
		b.ParentBOMID = &parentID.Int64
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		b.DeletedAt = &t
	}
	return b, nil
}

func ListBOMLines(db *sql.DB, filter BOMFilter) ([]BOMLine, error) {
	query := `
		SELECT id, product_id, component_product_id, component_unit_id,
			component_qty, waste_percentage, parent_bom_id, created_at, updated_at, deleted_at
		FROM bills_of_materials
		WHERE deleted_at IS NULL`
	args := []any{}

	if filter.ProductID != nil {
		query += " AND product_id = ?"
		args = append(args, *filter.ProductID)
	}
	query += " ORDER BY product_id, id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bomLines := []BOMLine{}
	for rows.Next() {
		var b BOMLine
		var parentID sql.NullInt64
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&b.ID,
			&b.ProductID,
			&b.ComponentProductID,
			&b.ComponentUnitID,
			&b.ComponentQty,
			&b.WastePercentage,
			&parentID,
			&b.CreatedAt,
			&b.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if parentID.Valid {
			b.ParentBOMID = &parentID.Int64
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			b.DeletedAt = &t
		}
		bomLines = append(bomLines, b)
	}

	return bomLines, rows.Err()
}

func UpdateBOMLine(db *sql.DB, b BOMLine) error {
	result, err := db.Exec(`
		UPDATE bills_of_materials SET
			product_id = ?, component_product_id = ?, component_unit_id = ?,
			component_qty = ?, waste_percentage = ?, parent_bom_id = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		b.ProductID,
		b.ComponentProductID,
		b.ComponentUnitID,
		b.ComponentQty,
		b.WastePercentage,
		nullableInt64Value(b.ParentBOMID),
		b.ID,
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

func DeleteBOMLine(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE bills_of_materials SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
