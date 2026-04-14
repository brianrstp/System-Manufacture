package database

import (
	"database/sql"
	"time"
)

type StockMovement struct {
	ID            int64      `json:"id"`
	InventoryID   *int64     `json:"inventoryId,omitempty"`
	ProductID     int64      `json:"productId"`
	WarehouseID   int64      `json:"warehouseId"`
	MovementType  string     `json:"movementType"`
	Quantity      float64    `json:"quantity"`
	UnitID        *int64     `json:"unitId,omitempty"`
	ReferenceType string     `json:"referenceType"`
	ReferenceID   string     `json:"referenceId"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

type StockMovementFilter struct {
	ProductID    *int64
	WarehouseID  *int64
	MovementType string
}

func CreateStockMovement(db *sql.DB, m StockMovement) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO stock_movements (
			inventory_id, product_id, warehouse_id, movement_type,
			quantity, unit_id, reference_type, reference_id, notes,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		nullableInt64Value(m.InventoryID),
		m.ProductID,
		m.WarehouseID,
		m.MovementType,
		m.Quantity,
		nullableInt64Value(m.UnitID),
		m.ReferenceType,
		m.ReferenceID,
		m.Notes,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetStockMovementByID(db *sql.DB, id int64) (StockMovement, error) {
	var m StockMovement
	var inventoryID sql.NullInt64
	var unitID sql.NullInt64
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, inventory_id, product_id, warehouse_id, movement_type,
			quantity, unit_id, reference_type, reference_id, notes,
			created_at, updated_at, deleted_at
		FROM stock_movements
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&m.ID,
		&inventoryID,
		&m.ProductID,
		&m.WarehouseID,
		&m.MovementType,
		&m.Quantity,
		&unitID,
		&m.ReferenceType,
		&m.ReferenceID,
		&m.Notes,
		&m.CreatedAt,
		&m.UpdatedAt,
		&deletedAt,
	); err != nil {
		return StockMovement{}, err
	}
	if inventoryID.Valid {
		m.InventoryID = &inventoryID.Int64
	}
	if unitID.Valid {
		m.UnitID = &unitID.Int64
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		m.DeletedAt = &t
	}
	return m, nil
}

func ListStockMovements(db *sql.DB, filter StockMovementFilter) ([]StockMovement, error) {
	query := `
		SELECT id, inventory_id, product_id, warehouse_id, movement_type,
			quantity, unit_id, reference_type, reference_id, notes,
			created_at, updated_at, deleted_at
		FROM stock_movements
		WHERE deleted_at IS NULL`
	args := []any{}

	if filter.ProductID != nil {
		query += " AND product_id = ?"
		args = append(args, *filter.ProductID)
	}
	if filter.WarehouseID != nil {
		query += " AND warehouse_id = ?"
		args = append(args, *filter.WarehouseID)
	}
	if filter.MovementType != "" {
		query += " AND movement_type = ?"
		args = append(args, filter.MovementType)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []StockMovement{}
	for rows.Next() {
		var m StockMovement
		var inventoryID sql.NullInt64
		var unitID sql.NullInt64
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&m.ID,
			&inventoryID,
			&m.ProductID,
			&m.WarehouseID,
			&m.MovementType,
			&m.Quantity,
			&unitID,
			&m.ReferenceType,
			&m.ReferenceID,
			&m.Notes,
			&m.CreatedAt,
			&m.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if inventoryID.Valid {
			m.InventoryID = &inventoryID.Int64
		}
		if unitID.Valid {
			m.UnitID = &unitID.Int64
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			m.DeletedAt = &t
		}
		movements = append(movements, m)
	}

	return movements, rows.Err()
}

func UpdateStockMovement(db *sql.DB, m StockMovement) error {
	result, err := db.Exec(`
		UPDATE stock_movements SET
			inventory_id = ?, product_id = ?, warehouse_id = ?, movement_type = ?,
			quantity = ?, unit_id = ?, reference_type = ?, reference_id = ?, notes = ?,
			updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		nullableInt64Value(m.InventoryID),
		m.ProductID,
		m.WarehouseID,
		m.MovementType,
		m.Quantity,
		nullableInt64Value(m.UnitID),
		m.ReferenceType,
		m.ReferenceID,
		m.Notes,
		m.ID,
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

func DeleteStockMovement(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE stock_movements SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
