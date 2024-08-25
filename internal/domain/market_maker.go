package domain

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

const (
	Taker = "taker"
	Maker = "maker"
)

type RecordData struct {
	Date          time.Time       `json:"date"`
	ActualCost    decimal.Decimal `json:"actual_cost"`
	EstimatedCost decimal.Decimal `json:"estimated_cost"`

	DexName              string          `json:"dex_name"`
	SwappedSymbol0       symbol.Symbol   `json:"swapped_symbol_0"`
	SwappedSymbol1       symbol.Symbol   `json:"swapped_symbol_1"`
	SwappedSymbol0Amount decimal.Decimal `json:"swapped_symbol_0_amount"`
	SwappedSymbol1Amount decimal.Decimal `json:"swapped_symbol_1_amount"`

	ExchangeName        string          `json:"exchange_name"`
	TradedSymbol0       symbol.Symbol   `json:"traded_symbol_0"`
	TradedSymbol1       symbol.Symbol   `json:"traded_symbol_1"`
	TradedSymbol0Amount decimal.Decimal `json:"traded_symbol_0_amount"`
	TradedSymbol1Amount decimal.Decimal `json:"traded_symbol_1_amount"`

	WalletToken0Balance         decimal.Decimal `json:"wallet_token_0_balance"`
	WalletToken1Balance         decimal.Decimal `json:"wallet_token_1_balance"`
	NobitexWalletMarket0Balance decimal.Decimal `json:"nobitex_wallet_market_0_balance"`
	NobitexWalletMarket1Balance decimal.Decimal `json:"nobitex_wallet_market_1_balance"`
	Profit                      decimal.Decimal `json:"profit"`

	CoinDexPrice           decimal.Decimal `json:"coin_dex_price"`
	CoinNobitexPrice       decimal.Decimal `json:"coin_nobitex_price"`
	UsdtPrice              decimal.Decimal `json:"usdt_price"`
	TxHash                 string          `json:"tx_hash"`
	DexPriceEstimation     decimal.Decimal `json:"dex_price_estimation"`
	NobitexPriceEstimation decimal.Decimal `json:"nobitex_price_estimation"`
	StrategyType           string          `json:"strategy_type"`
	State                  string          `json:"state"`
}
