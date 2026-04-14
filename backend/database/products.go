package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID              int64      `json:"id"`
	SKU             string     `json:"sku"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	CategoryID      *int64     `json:"categoryId,omitempty"`
	UnitID          *int64     `json:"unitId,omitempty"`
	ProductType     string     `json:"productType"`
	StandardPrice   float64    `json:"standardPrice"`
	CostPrice       float64    `json:"costPrice"`
	LeadTimeDays    int        `json:"leadTimeDays"`
	MinOrderQty     int        `json:"minOrderQty"`
	ReorderPoint    int        `json:"reorderPoint"`
	LifecycleStatus string     `json:"lifecycleStatus"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}

type ProductFilter struct {
	Search          string
	CategoryID      *int64
	ProductType     string
	LifecycleStatus string
}

func EnsureManufacturingSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			parent_id BIGINT UNSIGNED NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (parent_id),
			INDEX (status),
			FOREIGN KEY (parent_id) REFERENCES categories(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS units (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			factor DECIMAL(18,6) NOT NULL DEFAULT 1.000000,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS products (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			sku VARCHAR(100) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			category_id BIGINT UNSIGNED NULL,
			unit_id BIGINT UNSIGNED NULL,
			product_type VARCHAR(50) NOT NULL DEFAULT 'raw_material',
			standard_price DECIMAL(18,4) NOT NULL DEFAULT 0.00,
			cost_price DECIMAL(18,4) NOT NULL DEFAULT 0.00,
			lead_time_days INT UNSIGNED NOT NULL DEFAULT 0,
			min_order_qty INT UNSIGNED NOT NULL DEFAULT 1,
			reorder_point INT UNSIGNED NOT NULL DEFAULT 0,
			lifecycle_status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (product_type),
			INDEX (lifecycle_status),
			INDEX (category_id),
			INDEX (unit_id),
			FOREIGN KEY (category_id) REFERENCES categories(id),
			FOREIGN KEY (unit_id) REFERENCES units(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS product_images (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			product_id BIGINT UNSIGNED NOT NULL,
			url VARCHAR(1024) NOT NULL,
			alt_text VARCHAR(512),
			sort_order INT UNSIGNED NOT NULL DEFAULT 0,
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (product_id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS product_variants (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			product_id BIGINT UNSIGNED NOT NULL,
			sku VARCHAR(100),
			attributes JSON NULL,
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (product_id),
			INDEX (sku),
			FOREIGN KEY (product_id) REFERENCES products(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS bills_of_materials (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			product_id BIGINT UNSIGNED NOT NULL,
			component_product_id BIGINT UNSIGNED NOT NULL,
			component_unit_id BIGINT UNSIGNED NOT NULL,
			component_qty DECIMAL(18,4) NOT NULL DEFAULT 0.00,
			waste_percentage DECIMAL(5,2) NOT NULL DEFAULT 0.00,
			parent_bom_id BIGINT UNSIGNED NULL,
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (product_id, component_product_id),
			INDEX (parent_bom_id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (component_product_id) REFERENCES products(id),
			FOREIGN KEY (component_unit_id) REFERENCES units(id),
			FOREIGN KEY (parent_bom_id) REFERENCES bills_of_materials(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS warehouses (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(100) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			location VARCHAR(512),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS inventory (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			product_id BIGINT UNSIGNED NOT NULL,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			qty_on_hand DECIMAL(18,4) NOT NULL DEFAULT 0.0000,
			qty_reserved DECIMAL(18,4) NOT NULL DEFAULT 0.0000,
			qty_available DECIMAL(18,4) AS (qty_on_hand - qty_reserved) STORED,
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			UNIQUE KEY idx_product_warehouse (product_id, warehouse_id),
			INDEX (product_id),
			INDEX (warehouse_id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS customers (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NULL,
			phone VARCHAR(100) NULL,
			address TEXT NULL,
			password_hash VARCHAR(255) NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`ALTER TABLE customers ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) NULL;`,

		`CREATE TABLE IF NOT EXISTS orders (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			order_number VARCHAR(100) NOT NULL UNIQUE,
			customer_id BIGINT UNSIGNED NOT NULL,
			product VARCHAR(255) NOT NULL,
			order_date DATETIME(3) NOT NULL,
			amount DECIMAL(18,4) NOT NULL DEFAULT 0.00,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (order_date),
			INDEX (status),
			INDEX (customer_id),
			FOREIGN KEY (customer_id) REFERENCES customers(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS production_jobs (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			job_code VARCHAR(100) NOT NULL UNIQUE,
			product VARCHAR(255) NOT NULL,
			start_date DATETIME(3) NOT NULL,
			duration_days INT UNSIGNED NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (start_date),
			INDEX (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS stock_movements (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			inventory_id BIGINT UNSIGNED NULL,
			product_id BIGINT UNSIGNED NOT NULL,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			movement_type VARCHAR(50) NOT NULL,
			quantity DECIMAL(18,4) NOT NULL DEFAULT 0.0000,
			unit_id BIGINT UNSIGNED NULL,
			reference_type VARCHAR(100) NULL,
			reference_id VARCHAR(255) NULL,
			notes TEXT,
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
			deleted_at DATETIME(3) NULL DEFAULT NULL,
			INDEX (inventory_id),
			INDEX (product_id, warehouse_id),
			INDEX (movement_type),
			FOREIGN KEY (inventory_id) REFERENCES inventory(id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
			FOREIGN KEY (unit_id) REFERENCES units(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, sqlStmt := range statements {
		if _, err := db.Exec(sqlStmt); err != nil {
			return fmt.Errorf("gagal membuat schema: %w", err)
		}
	}

	return nil
}

func CreateProduct(db *sql.DB, p Product) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO products (
			sku, name, description, category_id, unit_id, product_type,
			standard_price, cost_price, lead_time_days, min_order_qty,
			reorder_point, lifecycle_status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		p.SKU,
		p.Name,
		p.Description,
		nullableInt64Value(p.CategoryID),
		nullableInt64Value(p.UnitID),
		p.ProductType,
		p.StandardPrice,
		p.CostPrice,
		p.LeadTimeDays,
		p.MinOrderQty,
		p.ReorderPoint,
		p.LifecycleStatus,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetProductByID(db *sql.DB, id int64) (Product, error) {
	var p Product
	var categoryID sql.NullInt64
	var unitID sql.NullInt64
	var deletedAt sql.NullTime

	row := db.QueryRow(`
		SELECT id, sku, name, description, category_id, unit_id, product_type,
			standard_price, cost_price, lead_time_days, min_order_qty,
			reorder_point, lifecycle_status, created_at, updated_at, deleted_at
		FROM products
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&categoryID,
		&unitID,
		&p.ProductType,
		&p.StandardPrice,
		&p.CostPrice,
		&p.LeadTimeDays,
		&p.MinOrderQty,
		&p.ReorderPoint,
		&p.LifecycleStatus,
		&p.CreatedAt,
		&p.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Product{}, err
	}
	if categoryID.Valid {
		p.CategoryID = &categoryID.Int64
	}
	if unitID.Valid {
		p.UnitID = &unitID.Int64
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		p.DeletedAt = &t
	}

	return p, nil
}

func ListProducts(db *sql.DB, filter ProductFilter) ([]Product, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT id, sku, name, description, category_id, unit_id, product_type,
			standard_price, cost_price, lead_time_days, min_order_qty,
			reorder_point, lifecycle_status, created_at, updated_at, deleted_at
		FROM products
		WHERE deleted_at IS NULL`)

	args := []any{}

	if filter.Search != "" {
		query.WriteString(" AND (sku LIKE ? OR name LIKE ? OR description LIKE ?)")
		term := "%" + filter.Search + "%"
		args = append(args, term, term, term)
	}
	if filter.CategoryID != nil {
		query.WriteString(" AND category_id = ?")
		args = append(args, *filter.CategoryID)
	}
	if filter.ProductType != "" {
		query.WriteString(" AND product_type = ?")
		args = append(args, filter.ProductType)
	}
	if filter.LifecycleStatus != "" {
		query.WriteString(" AND lifecycle_status = ?")
		args = append(args, filter.LifecycleStatus)
	}
	query.WriteString(" ORDER BY updated_at DESC")

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		var categoryID sql.NullInt64
		var unitID sql.NullInt64
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Description,
			&categoryID,
			&unitID,
			&p.ProductType,
			&p.StandardPrice,
			&p.CostPrice,
			&p.LeadTimeDays,
			&p.MinOrderQty,
			&p.ReorderPoint,
			&p.LifecycleStatus,
			&p.CreatedAt,
			&p.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			p.CategoryID = &categoryID.Int64
		}
		if unitID.Valid {
			p.UnitID = &unitID.Int64
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			p.DeletedAt = &t
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func UpdateProduct(db *sql.DB, p Product) error {
	result, err := db.Exec(`
		UPDATE products SET
			sku = ?, name = ?, description = ?, category_id = ?, unit_id = ?,
			product_type = ?, standard_price = ?, cost_price = ?, lead_time_days = ?,
			min_order_qty = ?, reorder_point = ?, lifecycle_status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		p.SKU,
		p.Name,
		p.Description,
		nullableInt64Value(p.CategoryID),
		nullableInt64Value(p.UnitID),
		p.ProductType,
		p.StandardPrice,
		p.CostPrice,
		p.LeadTimeDays,
		p.MinOrderQty,
		p.ReorderPoint,
		p.LifecycleStatus,
		p.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteProduct(db *sql.DB, id int64) error {
	result, err := db.Exec("UPDATE products SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func nullableInt64Value(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}
