package store

import (
	"context"
	"time"

	"github.com/zarbanio/market-maker-keeper/internal/domain"
)

func (p postgres) CreateCycle(ctx context.Context, start time.Time, status domain.CycleStatus) error {
	query := `
		INSERT INTO cycles ("start", status)
		VALUES ($1, $2);
		`

	_, err := p.conn.Exec(ctx, query, start, status)
	if err != nil {
		return err
	}
	return nil
}

func (p postgres) UpdateCycle(ctx context.Context, id int64, end time.Time, status domain.CycleStatus) error {
	query := `
		UPDATE cycles
		SET "end" = $2, status = $3
		WHERE id = $1;
		`
	_, err := p.conn.Exec(ctx, query, id, end, status)
	if err != nil {
		return err
	}
	return nil
}

func (p postgres) GetLastCycleId(ctx context.Context) (int64, error) {
	query := `
		SELECT id
		FROM cycles
		ORDER BY id DESC
		LIMIT 1;
		`

	var lastCycleId int64
	err := p.conn.QueryRow(ctx, query).Scan(&lastCycleId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return lastCycleId, nil
}
