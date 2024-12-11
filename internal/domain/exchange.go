package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
)

type Exchange interface {
	OrderBook(src, dst symbol.Symbol) (orderbook.OrderBook, error)
	ExchangeRate(src, dst symbol.Symbol) (price decimal.Decimal, err error)
	CancelOrder(order order.Order) error
	Balances() ([]Balance, error)
	SetBalance(symbol symbol.Symbol, balance decimal.Decimal)
	PlaceOrder(order order.Order) (string, time.Time, error)
	OrderStatus(ctx context.Context, id string) (order.Order, error)
	Orders(side order.Side, src, dst symbol.Symbol, state order.State) ([]order.Order, error)
	UpdateOrder(id string, state order.State) (order.State, error)
	RecentTrades(src, dst symbol.Symbol) ([]trade.Trade, error)
	Fees(tradeType trade.Type) decimal.Decimal
	Decimals(s symbol.Symbol) int32
	MinimumOrderToman() decimal.Decimal
}
