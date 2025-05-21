package repository

import (
	"cashback-serv/models"
	"database/sql"
	"time"
)

type SourceRepository struct {
	db *sql.DB
}

func NewSourceRepository(db *sql.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) CreateSource(source *models.Source) error {
	query := `
		INSERT INTO sources (
			host_ip,
			slug,
			created_at,
			updated_at
		) VALUES (
			$host_ip$,
			$slug$,
			$created_at$,
			$updated_at$
		) RETURNING id`

	now := time.Now()
	args := map[string]interface{}{
		"$host_ip$":    source.HostIP,
		"$slug$":       source.Slug,
		"$created_at$": now,
		"$updated_at$": now,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	return r.db.QueryRow(namedQuery, namedArgs...).Scan(&source.ID)
}

func (r *SourceRepository) GetSourceBySlug(slug string) (*models.Source, error) {
	query := `
		SELECT 
			id,
			host_ip,
			slug,
			created_at,
			updated_at,
			deleted_at
		FROM sources
		WHERE slug = $slug$ 
		AND deleted_at IS NULL`

	args := map[string]interface{}{
		"$slug$": slug,
	}

	namedQuery, namedArgs := buildNamedQuery(query, args)
	source := &models.Source{}
	err := r.db.QueryRow(namedQuery, namedArgs...).Scan(
		&source.ID,
		&source.HostIP,
		&source.Slug,
		&source.CreatedAt,
		&source.UpdatedAt,
		&source.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return source, err
}
