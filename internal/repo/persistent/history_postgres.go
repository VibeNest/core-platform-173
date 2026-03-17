package persistent

import (
	"context"
	"fmt"

	"gitverse.ru/apavlov-systems/core-platform/internal/entity"
	"gitverse.ru/apavlov-systems/core-platform/internal/repo/postgres"
)

type HistoryRepo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg: pg}
}

func (r *HistoryRepo) StoreHistory(ctx context.Context, t entity.TranslationHistory) error {
	sql := `INSERT INTO history (id, source, destination, original, translation, created_at) 
	        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pg.Pool.Exec(ctx, sql, t.ID, t.Source, t.Destination, t.Original, t.Translation, t.CreatedAt)
	if err != nil {
		return fmt.Errorf("HistoryRepo - StoreHistory - Exec: %w", err)
	}

	return nil
}

func (r *HistoryRepo) GetHistory(ctx context.Context) ([]entity.TranslationHistory, error) {
	sql := `SELECT id, source, destination, original, translation, created_at 
	        FROM history ORDER BY created_at DESC LIMIT 50`

	rows, err := r.pg.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("HistoryRepo - GetHistory - Query: %w", err)
	}
	defer rows.Close()

	var histories []entity.TranslationHistory

	for rows.Next() {
		var h entity.TranslationHistory
		err = rows.Scan(&h.ID, &h.Source, &h.Destination, &h.Original, &h.Translation, &h.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("HistoryRepo - GetHistory - Scan: %w", err)
		}
		histories = append(histories, h)
	}

	return histories, nil
}
