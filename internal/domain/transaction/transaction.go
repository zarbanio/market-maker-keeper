package transaction

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zarbanio/market-maker-keeper/internal/domain/blockchain"
)

type (
	Transaction struct {
		Id                int64                  `json:"id"`
		TxHash            common.Hash            `json:"tx_hash"`
		FromAddress       common.Address         `json:"from_address"`
		BlockchainId      *blockchain.Blockchain `json:"blockchain_id"`
		Timestamp         time.Time              `json:"timestamp"`
		BlockNumber       int64                  `json:"block_number"`
		ToAddress         common.Address         `json:"to_address"`
		Value             string                 `json:"value"`
		GasPrice          string                 `json:"gas_fee"`
		GasUsage          uint64                 `json:"gas_usage"`
		TransactionStatus State                  `json:"transaction_status"`
		TransactionData   *types.Transaction     `json:"transaction_data"`
	}
	UpdatedFields struct {
		BlockNumber int64
		GasUsed     uint64
		Timestamp   time.Time
		Status      State
	}
	State uint
)

const (
	Pending State = iota
	Success
	Failed
)

func (o State) String() string {
	return []string{"pending", "success", "failed"}[o]
}

func CastState(key string) State {
	var status State
	switch key {
	case "pending":
		status = Pending
	case "success":
		status = Success
	case "failed":
		status = Failed
	}
	return status
}

func CastFromReceiptStatus(s uint64) State {
	var status State
	switch s {
	case types.ReceiptStatusSuccessful:
		return Success
	case types.ReceiptStatusFailed:
		return Failed
	}
	return status
}
