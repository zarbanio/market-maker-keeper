package store

import (
	"context"
	"fmt"

	"github.com/zarbanio/market-maker-keeper/internal/domain/pair"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"

	"github.com/jackc/pgx/v4"
)

type pairModel struct {
	id         int64
	baseAsset  string
	quoteAsset string
}

func (p pairModel) toDomain() (error, *pair.Pair) {
	baseAsset, err := symbol.FromString(p.baseAsset)
	if err != nil {
		return err, nil
	}
	quoteAsset, err := symbol.FromString(p.quoteAsset)
	if err != nil {
		return err, nil
	}

	return nil, &pair.Pair{
		Id:         p.id,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}
}

func (p postgres) CreatePair(ctx context.Context, pair *pair.Pair) (int64, error) {

	// prepare insert statement
	stmt := `
        INSERT INTO pairs (base_asset, quote_asset)
		VALUES ($1, $2)
		RETURNING id
    `
	var id int64
	err := p.conn.QueryRow(ctx, stmt, pair.BaseAsset.String(), pair.QuoteAsset.String()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert pair: %v", err)
	}

	return id, nil
}

func (p postgres) CreatePairIfNotExist(ctx context.Context, pair *pair.Pair) (int64, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT id
		FROM pairs
		WHERE base_asset = $1 AND quote_asset = $2
	`, pair.BaseAsset.String(), pair.QuoteAsset.String())

	var id int64
	err := row.Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return p.CreatePair(ctx, pair)
		}
		return 0, fmt.Errorf("failed to get pair: %w", err)
	}
	return id, nil
}

func (p postgres) GetPairById(ctx context.Context, id int64) (*pair.Pair, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT *
		FROM pairs
		WHERE id = $1
	`, id)

	var pm pairModel
	err := row.Scan(&pm.id, &pm.baseAsset, &pm.quoteAsset)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrPairNotFound
		}
		return nil, fmt.Errorf("failed to get pair: %w", err)
	}
	err, normalizedPm := pm.toDomain()
	if err != nil {
		return nil, err
	}
	return normalizedPm, nil
}

func (p postgres) GetPairList(ctx context.Context) ([]pair.Pair, error) {
	const query = `SELECT * FROM pairs`

	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying pair: %w", err)
	}
	defer rows.Close()

	var PairList []pair.Pair
	for rows.Next() {
		var pm pairModel
		if err := rows.Scan(
			&pm.id,
			&pm.baseAsset,
			&pm.quoteAsset,
		); err != nil {
			return nil, fmt.Errorf("error scanning pair row: %w", err)
		}
		err, normalizedPm := pm.toDomain()
		if err != nil {
			return nil, err
		}
		PairList = append(PairList, *normalizedPm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over  pair rows: %w", err)
	}

	return PairList, nil
}
