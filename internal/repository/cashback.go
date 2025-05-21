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
		INSERT INTO "cashback" (
			cashback_amount,
			turon_user_id,
			created_at,
			updated_at
		) VALUES (
			$cashback_amount$,
			$turon_user_id$,
			$created_at$,
			$updated_at$
		) RETURNING id`

	now := time.Now()
	args := map[string]interface{}{
		"$cashback_amount$": cashback.CashbackAmount,
		"$turon_user_id$":   cashback.TuronUserID,
		"$created_at$":      now,
		"$updated_at$":      now,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	return r.db.QueryRow(namedQuery, namedArgs...).Scan(&cashback.ID)
}

func (r *CashbackRepository) CreateCashbackHistory(history *models.CashbackHistory) error {
	query := `
		INSERT INTO "cashback_history" (
			cashback_id,
			source_id,
			cashback_amount,
			host_ip,
			type,
			created_at,
			updated_at
		) VALUES (
			$cashback_id$,
			$source_id$,
			$cashback_amount$,
			$host_ip$,
			$type$,
			$created_at$,
			$updated_at$
		) RETURNING id`

	now := time.Now()
	args := map[string]interface{}{
		"$cashback_id$":     history.CashbackID,
		"$source_id$":       history.SourceID,
		"$cashback_amount$": history.CashbackAmount,
		"$host_ip$":         history.HostIP,
		"$type$":            history.Type,
		"$created_at$":      now,
		"$updated_at$":      now,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	return r.db.QueryRow(namedQuery, namedArgs...).Scan(&history.ID)
}

func (r *CashbackRepository) GetCashbackByUserID(turonUserID int64) (*models.Cashback, error) {
	query := `
		SELECT 
			id,
			cashback_amount,
			turon_user_id,
			created_at,
			updated_at,
			deleted_at
		FROM cashback
		WHERE turon_user_id = $turon_user_id$
		AND deleted_at IS NULL`

	args := map[string]interface{}{
		"$turon_user_id$": turonUserID,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	cashback := &models.Cashback{}
	err := r.db.QueryRow(namedQuery, namedArgs...).Scan(
		&cashback.ID,
		&cashback.CashbackAmount,
		&cashback.TuronUserID,
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
		UPDATE cashback
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

func (r *CashbackRepository) buildDateFilters(query string, args map[string]interface{}, fromDate, toDate string) string {
	if fromDate != "" {
		query += " AND ch.created_at >= $from_date$"
		args["$from_date$"] = fromDate
	}
	if toDate != "" {
		query += " AND ch.created_at <= $to_date$"
		args["$to_date$"] = toDate
	}
	return query
}

func (r *CashbackRepository) buildPagination(query string, args map[string]interface{}, pagination *models.Pagination) string {
	query += " ORDER BY ch.created_at DESC LIMIT $limit$ OFFSET $offset$"
	args["$limit$"] = pagination.Limit
	args["$offset$"] = pagination.Offset
	return query
}

func (r *CashbackRepository) GetCashbackHistoryByUserID(turonUserID int64, fromDate, toDate string, pagination *models.Pagination) ([]models.CashbackHistory, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM cashback_history ch
		JOIN cashback c ON c.id = ch.cashback_id
		WHERE c.turon_user_id = $turon_user_id$ 
		AND ch.deleted_at IS NULL`

	args := map[string]interface{}{
		"$turon_user_id$": turonUserID,
	}

	countQuery = r.buildDateFilters(countQuery, args, fromDate, toDate)

	namedCountQuery, namedCountArgs := buildNamedQuery(countQuery, args)
	var total int64
	err := r.db.QueryRow(namedCountQuery, namedCountArgs...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	pagination.ItemTotal = total
	pagination.PageTotal = (total + pagination.PageSize - 1) / pagination.PageSize

	query := `
		SELECT 
			ch.id,
			ch.cashback_id,
			s.slug as source_slug,
			ch.cashback_amount,
			ch.host_ip,
			ch.type,
			ch.created_at,
			ch.updated_at,
			ch.deleted_at
		FROM cashback_history ch
		JOIN cashback c ON c.id = ch.cashback_id
		LEFT JOIN sources s ON s.id = ch.source_id
		WHERE c.turon_user_id = $turon_user_id$ 
		AND ch.deleted_at IS NULL`

	query = r.buildDateFilters(query, args, fromDate, toDate)
	query = r.buildPagination(query, args, pagination)

	namedQuery, namedArgs := buildNamedQuery(query, args)
	rows, err := r.db.Query(namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cashback history: %w", err)
	}
	defer rows.Close()

	var history []models.CashbackHistory
	for rows.Next() {
		var h models.CashbackHistory
		var sourceSlug sql.NullString
		if err := rows.Scan(
			&h.ID,
			&h.CashbackID,
			&sourceSlug,
			&h.CashbackAmount,
			&h.HostIP,
			&h.Type,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cashback history row: %w", err)
		}
		if sourceSlug.Valid {
			h.SourceSlug = sourceSlug.String
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
