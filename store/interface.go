package store

import (
	"context"
	"time"

	"github.com/zarbanio/market-maker-keeper/internal/domain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/pair"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"
)

type IStore interface {
	MigrateTable
	ITransaction
	IOrder
	ITrade
	IPair
	IBlockPtr
	ICycle
	ILog
	IArbitrageOpportunity
}

type MigrateTable interface {
	Migrate(path string) error
}

type ITransaction interface {
	CreateTransaction(ctx context.Context, tx *types.Transaction, from common.Address) (int64, error)
	UpdateTransaction(ctx context.Context, transactionID int64, updatedFields transaction.UpdatedFields) error
	GetTransactionById(ctx context.Context, id int64) (*transaction.Transaction, error)
	GetTransactionIdByHash(ctx context.Context, transactionHash string) (int64, error)
}

type IPair interface {
	CreatePair(ctx context.Context, pair *pair.Pair) (int64, error)
	CreatePairIfNotExist(ctx context.Context, pair *pair.Pair) (int64, error)
	GetPairById(ctx context.Context, id int64) (*pair.Pair, error)
	GetPairList(ctx context.Context) ([]pair.Pair, error)
}

type IOrder interface {
	CreateNewOrder(ctx context.Context, order order.Order) (int64, error)
	GetOrderById(ctx context.Context, id int64) (*order.Order, error)
	GetOrderByOrderId(ctx context.Context, orderID int64) (*order.Order, error)
	GetOrders(ctx context.Context, page int, pageSize int) ([]order.Order, error)
	UpdateOrder(ctx context.Context, orderID int64, updatedFields order.UpdatedFields) error
}

type ITrade interface {
	CreateNewTrade(ctx context.Context, pairId int64, orderId int64, transactionId int64) (int64, error)
	GetTradeByID(ctx context.Context, id int64) (*trade.Trade, error)
	//GetTrades(ctx context.Context) ([]trade.Trade, error)
}

type IBlockPtr interface {
	ExistsBlockPtr(ctx context.Context) (bool, error)
	CreateBlockPtr(ctx context.Context, ptr uint64) (int64, error)
	GetBlockPtrById(ctx context.Context, id int64) (uint64, error)
	GetLastBlockPtr(ctx context.Context) (int64, uint64, error)
	UpdateBlockPtr(ctx context.Context, id int64, ptr uint64) (uint64, error)
	IncBlockPtr(ctx context.Context, id int64) (uint64, error)
	DeleteBlockPtr(ctx context.Context, id int64) error
}

type ICycle interface {
	CreateCycle(ctx context.Context, start time.Time, status domain.CycleStatus) error
	UpdateCycle(ctx context.Context, id int64, end time.Time, status domain.CycleStatus) error
	GetLastCycleId(ctx context.Context) (int64, error)
}

type ILog interface {
	CreateLog(ctx context.Context, b []byte) (int, error)
}

type IArbitrageOpportunity interface {
	CreateArbitrageOpporchunity(ctx context.Context, opportunity []byte) error
}
