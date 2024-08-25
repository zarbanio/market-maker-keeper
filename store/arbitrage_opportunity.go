package store

import (
	"context"
	"fmt"
	"time"
)

func (p postgres) CreateArbitrageOpporchunity(ctx context.Context, opportunity []byte) error {

	currentTimestamp := time.Now().Unix()
	stmt := `
        INSERT INTO arbitrage_opportunity (uuid, opportunity) 
		VALUES ($1, $2)
		RETURNING id
    `

	_, err := p.conn.Exec(context.Background(), stmt, currentTimestamp, opportunity)
	if err != nil {
		return fmt.Errorf("failed to insert arbitrage opportunity: %w", err)
	}
	return nil
}
