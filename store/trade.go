package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
)

type tradeModel struct {
	id            int64
	pairId        int64
	orderId       int64
	transactionId int64
}

func (p postgres) CreateNewTrade(ctx context.Context, pairId int64, orderId int64, transactionId int64) (int64, error) {
	// prepare insert statement
	stmt := `
        INSERT INTO trades (       
			pair_id,
			order_id,
			transaction_id
		) 
		VALUES ($1, $2, $3)
		RETURNING id
    `
	var id int64

	err := p.conn.QueryRow(context.Background(), stmt,
		pairId,
		orderId,
		transactionId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert trade: %w", err)
	}

	return id, nil
}

func (p *postgres) GetTradeByID(ctx context.Context, id int64) (*trade.Trade, error) {
	var t tradeModel

	err := p.conn.QueryRow(ctx, `
		SELECT id, pair_id, order_id, transaction_id
		FROM trades
		WHERE id = $1
	`, id).Scan(&t.id, &t.pairId, &t.orderId, &t.transactionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTradeNotFound
		}
		return nil, fmt.Errorf("failed : %s", err)
	}

	tx, err := p.GetTransactionById(ctx, t.transactionId)
	if err != nil {
		return nil, err
	}
	order, err := p.GetOrderById(ctx, t.orderId)
	if err != nil {
		return nil, err
	}
	pair, err := p.GetPairById(ctx, t.pairId)
	if err != nil {
		return nil, err
	}

	return &trade.Trade{
		Id:          t.id,
		Pair:        pair,
		Transaction: tx,
		Order:       order,
	}, nil
}

//func (p *postgres) GetTrades(ctx context.Context) ([]trade.Trade, error) {
//	rows, err := p.conn.Query(ctx, `
//        SELECT t.id, t.trading_pair_id, t.timestamp, t.price, t.quantity, t.side,
//			t.exchange_name, t.order_id, t.transaction_id, t.fee,
//			tx.tx_hash,
//			tx.blockchain_id,
//			tx.timestamp,
//			tx.block_number,
//			tx.transaction_status,
//			p.base_asset,
//			p.quote_asset
//        FROM trades t
//        LEFT JOIN transactions tx ON t.transaction_id = tx.id
//		LEFT JOIN pair p ON t.trading_pair_id = p.id
//    `)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var trades []trade.Trade
//	for rows.Next() {
//		t := &tradeModel{}
//		err = rows.Scan(
//			&t.id, &t.pair, &t.timestamp, &t.price, &t.quantity,
//			&t.side, &t.exchangeName, &t.order, &t.transactionId, &t.fee,
//			&t.transactionId.txHash, &t.transactionId.blockchainID,
//			&t.transactionId.timestamp, &t.transactionId.blockNumber, &t.transactionId.transactionStatus,
//			&t.pair.baseAsset, &t.pair.quoteAsset,
//		)
//		if err != nil {
//			return nil, err
//		}
//		normedValue, err := t.ToDomain()
//		if err != nil {
//			return nil, err
//		}
//		trades = append(trades, *normedValue)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//	return trades, nil
//}
