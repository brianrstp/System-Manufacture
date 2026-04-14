package database

import (
	"database/sql"
	"time"
)

type Inventory struct {
	ID           int64      `json:"id"`
	ProductID    int64      `json:"productId"`
	WarehouseID  int64      `json:"warehouseId"`
	QtyOnHand    float64    `json:"qtyOnHand"`
	QtyReserved  float64    `json:"qtyReserved"`
	QtyAvailable float64    `json:"qtyAvailable"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type InventoryFilter struct {
	ProductID   *int64
	WarehouseID *int64
}

func CreateInventory(db *sql.DB, inv Inventory) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO inventory (
			product_id, warehouse_id, qty_on_hand, qty_reserved, created_at, updated_at
		) VALUES (?, ?, ?, ?, NOW(3), NOW(3))`,
		inv.ProductID,
		inv.WarehouseID,
		inv.QtyOnHand,
		inv.QtyReserved,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetInventoryByID(db *sql.DB, id int64) (Inventory, error) {
	var inv Inventory
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, product_id, warehouse_id, qty_on_hand, qty_reserved, qty_available,
			created_at, updated_at, deleted_at
		FROM inventory
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&inv.ID,
		&inv.ProductID,
		&inv.WarehouseID,
		&inv.QtyOnHand,
		&inv.QtyReserved,
		&inv.QtyAvailable,
		&inv.CreatedAt,
		&inv.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Inventory{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		inv.DeletedAt = &t
	}
	return inv, nil
}

func ListInventory(db *sql.DB, filter InventoryFilter) ([]Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, qty_on_hand, qty_reserved, qty_available,
			created_at, updated_at, deleted_at
		FROM inventory
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
	query += " ORDER BY warehouse_id, product_id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventories := []Inventory{}
	for rows.Next() {
		var inv Inventory
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&inv.ID,
			&inv.ProductID,
			&inv.WarehouseID,
			&inv.QtyOnHand,
			&inv.QtyReserved,
			&inv.QtyAvailable,
			&inv.CreatedAt,
			&inv.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			inv.DeletedAt = &t
		}
		inventories = append(inventories, inv)
	}

	return inventories, rows.Err()
}

func UpdateInventory(db *sql.DB, inv Inventory) error {
	result, err := db.Exec(`
		UPDATE inventory SET qty_on_hand = ?, qty_reserved = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		inv.QtyOnHand,
		inv.QtyReserved,
		inv.ID,
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

func DeleteInventory(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE inventory SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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
