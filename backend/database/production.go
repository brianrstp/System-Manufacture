package database

import (
	"database/sql"
	"strings"
	"time"
)

type ProductionJob struct {
	ID           int64      `json:"id"`
	JobCode      string     `json:"jobCode"`
	Product      string     `json:"product"`
	StartDate    time.Time  `json:"startDate"`
	DurationDays int        `json:"durationDays"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type ProductionJobFilter struct {
	Search string
	Status string
}

func CreateProductionJob(db *sql.DB, job ProductionJob) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO production_jobs (
			job_code, product, start_date, duration_days, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, NOW(3), NOW(3))`,
		job.JobCode,
		job.Product,
		job.StartDate,
		job.DurationDays,
		job.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetProductionJobByID(db *sql.DB, id int64) (ProductionJob, error) {
	var job ProductionJob
	var deletedAt sql.NullTime
	row := db.QueryRow(`
		SELECT id, job_code, product, start_date, duration_days, status,
			created_at, updated_at, deleted_at
		FROM production_jobs
		WHERE id = ? AND deleted_at IS NULL`, id)

	if err := row.Scan(
		&job.ID,
		&job.JobCode,
		&job.Product,
		&job.StartDate,
		&job.DurationDays,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
		&deletedAt,
	); err != nil {
		return ProductionJob{}, err
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		job.DeletedAt = &t
	}
	return job, nil
}

func ListProductionJobs(db *sql.DB, filter ProductionJobFilter) ([]ProductionJob, error) {
	query := `
		SELECT id, job_code, product, start_date, duration_days, status,
			created_at, updated_at, deleted_at
		FROM production_jobs
		WHERE deleted_at IS NULL`
	args := []any{}

	if filter.Search != "" {
		query += ` AND (
			LOWER(job_code) LIKE ? OR
			LOWER(product) LIKE ? OR
			LOWER(status) LIKE ?
		)`
		pattern := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		args = append(args, pattern, pattern, pattern)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	query += " ORDER BY start_date DESC, id DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []ProductionJob{}
	for rows.Next() {
		var job ProductionJob
		var deletedAt sql.NullTime
		if err := rows.Scan(
			&job.ID,
			&job.JobCode,
			&job.Product,
			&job.StartDate,
			&job.DurationDays,
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			t := deletedAt.Time
			job.DeletedAt = &t
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func UpdateProductionJob(db *sql.DB, job ProductionJob) error {
	result, err := db.Exec(`
		UPDATE production_jobs SET job_code = ?, product = ?, start_date = ?, duration_days = ?, status = ?, updated_at = NOW(3)
		WHERE id = ? AND deleted_at IS NULL`,
		job.JobCode,
		job.Product,
		job.StartDate,
		job.DurationDays,
		job.Status,
		job.ID,
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

func DeleteProductionJob(db *sql.DB, id int64) error {
	result, err := db.Exec(`UPDATE production_jobs SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
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

func CountProductionJobs(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM production_jobs WHERE deleted_at IS NULL")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func CountProductionJobsByStatus(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query(`
		SELECT status, COUNT(*)
		FROM production_jobs
		WHERE deleted_at IS NULL
		GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}
