package trade

import (
	"fmt"

	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/pair"
	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"
)

const (
	Maker = "maker"
	Taker = "taker"
)

type (
	Type  string
	Trade struct {
		Id          int64                    `json:"id"`
		Pair        *pair.Pair               `json:"pair,omitempty"`
		Order       *order.Order             `json:"order"`
		Transaction *transaction.Transaction `json:"transaction,omitempty"`
	}
	Side uint
)

var sideMap = map[string]Side{
	"buy":  Buy,
	"sell": Sell,
}

const (
	Buy Side = iota
	Sell
)

func (s Side) String() string {
	return string([]string{"buy", "sell"}[s])
}

func SideFromString(s string) (Side, error) {
	side, ok := sideMap[s]
	if !ok {
		return 0, fmt.Errorf("invalid side %s", s)
	}
	return side, nil
}
