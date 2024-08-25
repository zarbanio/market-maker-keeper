package nobitex

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

func NewDecimalFromString(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}

func toNobitex(i interface{}) string {
	switch i.(type) {
	case symbol.Symbol:
		return strings.ToLower(i.(symbol.Symbol).String())
	case order.Side:
		return []string{"buy", "sell"}[i.(order.Side)]
	case order.Execution:
		return []string{"market", "limit", "stop_market", "stop_limit"}[i.(order.Execution)]
	case order.State:
		return []string{"open", "done", "done", "closed"}[i.(order.State)]
	}
	return ""
}
