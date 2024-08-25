package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"
)

type TransactionModel struct {
	id                int64
	txHash            string
	fromAddress       string
	timestamp         time.Time
	blockNumber       int64
	toAddress         string
	value             decimal.Decimal
	gasPrice          decimal.Decimal
	gasUsage          uint64
	transactionStatus string
	transactionData   json.RawMessage
}

func (t *TransactionModel) toDomain() (*transaction.Transaction, error) {
	var transactionData types.Transaction
	err := json.Unmarshal(t.transactionData, &transactionData)
	if err != nil {
		return nil, fmt.Errorf("to domain transaction failed: %v", err)
	}
	return &transaction.Transaction{
		Id:                t.id,
		TxHash:            common.HexToHash(t.txHash),
		FromAddress:       common.HexToAddress(t.fromAddress),
		Timestamp:         t.timestamp,
		BlockNumber:       t.blockNumber,
		ToAddress:         common.HexToAddress(t.toAddress),
		Value:             t.value.String(),
		GasPrice:          t.gasPrice.String(),
		GasUsage:          t.gasUsage,
		TransactionStatus: transaction.CastState(t.transactionStatus),
		TransactionData:   &transactionData,
	}, nil
}
func (p postgres) CreateTransaction(ctx context.Context, tx *types.Transaction, from common.Address) (int64, error) {
	txData, err := tx.MarshalJSON()
	if err != nil {
		return 0, err
	}
	stmt := `
        INSERT INTO transactions (
            tx_hash,
            from_address,
            to_address,
            value,
            gas_price,
            transaction_status,
            transaction_data
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
    `
	var id int64

	err = p.conn.QueryRow(context.Background(), stmt,
		tx.Hash().Hex(),
		from.String(),
		tx.To().Hex(),
		tx.Value().String(),
		tx.GasPrice().String(),
		"pending",
		string(txData),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %w", err)
	}

	return id, nil

}

func (p postgres) UpdateTransaction(ctx context.Context, transactionId int64, updatedFields transaction.UpdatedFields) error {
	// Prepare update statement
	stmt := `
		UPDATE transactions
		SET
			block_number = COALESCE($2, block_number),
			gas_usage = COALESCE($3, gas_usage),
			timestamp = COALESCE($4, timestamp),
			transaction_status = COALESCE($5, transaction_status)
		WHERE id = $1
	`

	_, err := p.conn.Exec(
		ctx,
		stmt,
		transactionId,
		updatedFields.BlockNumber,
		updatedFields.GasUsed,
		updatedFields.Timestamp.UTC(),
		updatedFields.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (p postgres) GetTransactionById(ctx context.Context, id int64) (*transaction.Transaction, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT *
		FROM transactions
		WHERE id = $1
	`, id)

	t := &TransactionModel{}
	var (
		blockNumber sql.NullInt64
		gasUsage    sql.NullInt64
		timestamp   sql.NullTime
	)
	err := row.Scan(
		&t.id,
		&t.txHash,
		&t.fromAddress,
		&timestamp,
		&blockNumber,
		&t.toAddress,
		&t.value,
		&t.gasPrice,
		&gasUsage,
		&t.transactionStatus,
		&t.transactionData,
	)
	if blockNumber.Valid {
		t.blockNumber = blockNumber.Int64
	}
	if gasUsage.Valid {
		t.gasUsage = uint64(gasUsage.Int64)
	}
	if timestamp.Valid {
		t.timestamp = timestamp.Time
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return t.toDomain()
}

func (p postgres) GetTransactionIdByHash(ctx context.Context, transactionHash string) (int64, error) {
	row := p.conn.QueryRow(ctx, `
        SELECT id
        FROM transactions
        WHERE tx_hash = $1
    `, transactionHash)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
