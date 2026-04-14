package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Order struct {
	ID           int64      `json:"id"`
	OrderNumber  string     `json:"orderNumber"`
	CustomerID   int64      `json:"customerId"`
	CustomerName string     `json:"customerName,omitempty"`
	Product      string     `json:"product"`
	OrderDate    time.Time  `json:"orderDate"`
	Amount       float64    `json:"amount"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type OrderFilter struct {
	Search     string
	Status     string
	CustomerID int64
	Limit      int
}

func CreateOrder(db *sql.DB, o Order) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO orders (
			order_number, customer_id, product, order_date, amount, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		o.OrderNumber,
		o.CustomerID,
		o.Product,
		o.OrderDate,
		o.Amount,
		o.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetOrderByID(db *sql.DB, id int64) (Order, error) {
	var o Order
	var deletedAt sql.NullTime
	row := db.QueryRow(`
		SELECT orders.id, orders.order_number, orders.customer_id, customers.name,
			orders.product, orders.order_date, orders.amount, orders.status,
			orders.created_at, orders.updated_at, orders.deleted_at
		FROM orders
		LEFT JOIN customers ON customers.id = orders.customer_id
		WHERE orders.id = ? AND orders.deleted_at IS NULL`, id)

	if err := row.Scan(
		&o.ID,
		&o.OrderNumber,
		&o.CustomerID,
		&o.CustomerName,
		&o.Product,
		&o.OrderDate,
		&o.Amount,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
		&deletedAt,
	); err != nil {
		return Order{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		o.DeletedAt = &t
	}
	return o, nil
}

func ListOrders(db *sql.DB, filter OrderFilter) ([]Order, error) {
	query := `
		SELECT orders.id, orders.order_number, orders.customer_id, customers.name,
			orders.product, orders.order_date, orders.amount, orders.status,
			orders.created_at, orders.updated_at, orders.deleted_at
		FROM orders
		LEFT JOIN customers ON customers.id = orders.customer_id
		WHERE orders.deleted_at IS NULL`
	args := []any{}

	if filter.Search != "" {
		query += ` AND (
			LOWER(order_number) LIKE ? OR
			LOWER(customers.name) LIKE ? OR
			LOWER(product) LIKE ? OR
			LOWER(status) LIKE ?
		)`
		pattern := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}
	if filter.CustomerID > 0 {
		query += " AND orders.customer_id = ?"
		args = append(args, filter.CustomerID)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	query += " ORDER BY order_date DESC, id DESC"
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var o Order
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&o.ID,
			&o.OrderNumber,
			&o.CustomerID,
			&o.CustomerName,
			&o.Product,
			&o.OrderDate,
			&o.Amount,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			o.DeletedAt = &t
		}
		orders = append(orders, o)
	}

	return orders, rows.Err()
}

func UpdateOrder(db *sql.DB, o Order) error {
	result, err := db.Exec(`
		UPDATE orders SET order_number = ?, customer_id = ?, product = ?, order_date = ?, amount = ?, status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		o.OrderNumber,
		o.CustomerID,
		o.Product,
		o.OrderDate,
		o.Amount,
		o.Status,
		o.ID,
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

func DeleteOrder(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE orders SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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

func CountOrders(db *sql.DB, status *string) (int, error) {
	query := "SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL"
	args := []any{}
	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}
	row := db.QueryRow(query, args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func SumOrderAmounts(db *sql.DB) (float64, error) {
	row := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM orders WHERE deleted_at IS NULL")
	var sum float64
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}

func SumOrderAmountsByMonth(db *sql.DB, months int) ([]float64, error) {
	values := make([]float64, months)
	query := `
		SELECT YEAR(order_date) AS year, MONTH(order_date) AS month, COALESCE(SUM(amount),0)
		FROM orders
		WHERE deleted_at IS NULL AND order_date >= DATE_SUB(CURDATE(), INTERVAL ? MONTH)
		GROUP BY year, month
		ORDER BY year, month`
	rows, err := db.Query(query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthMap := map[string]float64{}
	for rows.Next() {
		var year int
		var month int
		var sum float64
		if err := rows.Scan(&year, &month, &sum); err != nil {
			return nil, err
		}
		monthMap[fmt.Sprintf("%d-%02d", year, month)] = sum
	}

	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		key := fmt.Sprintf("%d-%02d", month.Year(), int(month.Month()))
		values[months-1-i] = monthMap[key]
	}

	return values, rows.Err()
}
