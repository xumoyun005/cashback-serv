package repository

import (
	"cashback-serv/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type CashbackRepository struct {
	db *sql.DB
}

func NewCashbackRepository(db *sql.DB) *CashbackRepository {
	return &CashbackRepository{db: db}
}

func (r *CashbackRepository) CreateCashback(cashback *models.Cashback) error {
	query := `
		INSERT INTO "cashbacks" (
			cashback_amount,
			turon_user_id,
			cinerama_user_id,
			created_at,
			updated_at
		) VALUES (
			$cashback_amount$,
			$turon_user_id$,
			$cinerama_user_id$,
			$created_at$,
			$updated_at$
		) RETURNING id`

	now := time.Now()
	args := map[string]interface{}{
		"$cashback_amount$":  cashback.CashbackAmount,
		"$turon_user_id$":    cashback.TuronUserID,
		"$cinerama_user_id$": cashback.CineramaUserID,
		"$created_at$":       now,
		"$updated_at$":       now,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	return r.db.QueryRow(namedQuery, namedArgs...).Scan(&cashback.ID)
}

func (r *CashbackRepository) CreateCashbackHistory(history *models.CashbackHistory) error {
	query := `
		INSERT INTO "cashback_histories" (
			cashback_id,
			cashback_amount,
			host_ip,
			device,
			type,
			created_at,
			updated_at
		) VALUES (
			$cashback_id$,
			$cashback_amount$,
			$host_ip$,
			$device$,
			$type$,
			$created_at$,
			$updated_at$
		) RETURNING id`

	now := time.Now()
	args := map[string]interface{}{
		"$cashback_id$":     history.CashbackID,
		"$cashback_amount$": history.CashbackAmount,
		"$host_ip$":         history.HostIP,
		"$device$":          history.Device,
		"$type$":            history.Type,
		"$created_at$":      now,
		"$updated_at$":      now,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	return r.db.QueryRow(namedQuery, namedArgs...).Scan(&history.ID)
}

func (r *CashbackRepository) GetCashbackByUserID(turonUserID int64, fromDate, toDate string) (*models.Cashback, error) {
	query := `
		SELECT 
			id,
			cashback_amount,
			turon_user_id,
			cinerama_user_id,
			created_at,
			updated_at,
			deleted_at
		FROM cashbacks
		WHERE turon_user_id = $turon_user_id$ 
		AND deleted_at IS NULL`

	args := map[string]interface{}{
		"$turon_user_id$": turonUserID,
	}

	if fromDate != "" {
		query += " AND created_at >= $from_date$"
		args["$from_date$"] = fromDate
	}

	if toDate != "" {
		query += " AND created_at <= $to_date$"
		args["$to_date$"] = toDate
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	cashback := &models.Cashback{}
	err := r.db.QueryRow(namedQuery, namedArgs...).Scan(
		&cashback.ID,
		&cashback.CashbackAmount,
		&cashback.TuronUserID,
		&cashback.CineramaUserID,
		&cashback.CreatedAt,
		&cashback.UpdatedAt,
		&cashback.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cashback, err
}

func (r *CashbackRepository) UpdateCashbackAmount(id int64, newAmount float64) error {
	query := `
		UPDATE cashbacks
		SET 
			cashback_amount = $cashback_amount$,
			updated_at = $updated_at$
		WHERE id = $id$ 
		AND deleted_at IS NULL`

	args := map[string]interface{}{
		"$cashback_amount$": newAmount,
		"$updated_at$":      time.Now(),
		"$id$":              id,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	_, err := r.db.Exec(namedQuery, namedArgs...)
	return err
}

func (r *CashbackRepository) GetCashbackHistoryByUserID(turonUserID int64, fromDate, toDate string, pagination *models.Pagination) ([]models.CashbackHistory, error) {
	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM cashback_histories ch
		JOIN cashbacks c ON c.id = ch.cashback_id
		WHERE c.turon_user_id = $turon_user_id$ 
		AND ch.deleted_at IS NULL`

	args := map[string]interface{}{
		"$turon_user_id$": turonUserID,
	}

	if fromDate != "" {
		countQuery += " AND ch.created_at >= $from_date$"
		args["$from_date$"] = fromDate
	}

	if toDate != "" {
		countQuery += " AND ch.created_at <= $to_date$"
		args["$to_date$"] = toDate
	}

	namedCountQuery, namedCountArgs := buildNamedQuery(countQuery, args)
	var total int64
	err := r.db.QueryRow(namedCountQuery, namedCountArgs...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pagination.ItemTotal = total
	pagination.PageTotal = (total + pagination.PageSize - 1) / pagination.PageSize

	// Get paginated data
	query := `
		SELECT 
			ch.id,
			ch.cashback_id,
			ch.cashback_amount,
			ch.host_ip,
			ch.device,
			ch.type,
			ch.created_at,
			ch.updated_at,
			ch.deleted_at
		FROM cashback_histories ch
		JOIN cashbacks c ON c.id = ch.cashback_id
		WHERE c.turon_user_id = $turon_user_id$ 
		AND ch.deleted_at IS NULL`

	if fromDate != "" {
		query += " AND ch.created_at >= $from_date$"
	}
	if toDate != "" {
		query += " AND ch.created_at <= $to_date$"
	}

	query += " ORDER BY ch.created_at DESC LIMIT $limit$ OFFSET $offset$"
	args["$limit$"] = pagination.Limit
	args["$offset$"] = pagination.Offset

	namedQuery, namedArgs := buildNamedQuery(query, args)
	rows, err := r.db.Query(namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cashback history: %w", err)
	}
	defer rows.Close()

	var history []models.CashbackHistory
	for rows.Next() {
		var h models.CashbackHistory
		if err := rows.Scan(
			&h.ID,
			&h.CashbackID,
			&h.CashbackAmount,
			&h.HostIP,
			&h.Device,
			&h.Type,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cashback history row: %w", err)
		}
		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cashback history rows: %w", err)
	}

	return history, nil
}

func buildNamedQuery(query string, args map[string]interface{}) (string, []interface{}) {
	var positionalArgs []interface{}
	position := 1

	for key, value := range args {
		query = replaceAll(query, key, fmt.Sprintf("$%d", position))
		positionalArgs = append(positionalArgs, value)
		position++
	}

	return query, positionalArgs
}

func replaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}
