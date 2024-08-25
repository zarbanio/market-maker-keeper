package order

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type (
	Order struct {
		Id              int64           `json:"id"`
		OrderId         int64           `json:"order_id"`
		Execution       Execution       `json:"execution"`
		Side            Side            `json:"side"`
		SrcCurrency     symbol.Symbol   `json:"srcCurrency"`
		DstCurrency     symbol.Symbol   `json:"dstCurrency"`
		Price           decimal.Decimal `json:"price"`
		Amount          decimal.Decimal `json:"amount"`
		TotalPrice      decimal.Decimal `json:"total_price"`
		TotalOrderPrice decimal.Decimal `json:"total_order_price"`
		StopPrice       decimal.Decimal `json:"stop_price"`
		Status          State           `json:"status"`
		Fee             decimal.Decimal `json:"fee"`
		FeeCurrency     symbol.Symbol   `json:"fee_currency"`
		Account         string          `json:"user"`
		CreatedAt       time.Time       `json:"created_at"`
		UnmatchedAmount decimal.Decimal `json:"unmatchedAmount"`
		MatchedAmount   decimal.Decimal `json:"matchedAmount"`
		Partial         bool            `json:"partial"`
	}
	UpdatedFields struct {
		MatchedAmount   *decimal.Decimal `json:"MatchedAmount,omitempty" bson:",omitempty"`
		UnmatchedAmount *decimal.Decimal `json:"UnmatchedAmount,omitempty" bson:",omitempty"`
		Status          *State           `json:"Status,omitempty" bson:",omitempty"`
		Fee             *decimal.Decimal `json:"Fee,omitempty" bson:",omitempty"`
		Price           *decimal.Decimal `json:"Price,omitempty" bson:",omitempty"`
		TotalPrice      *decimal.Decimal `json:"TotalPrice,omitempty" bson:"total_price"`
		TotalOrderPrice *decimal.Decimal `json:"TotalOrderPrice,omitempty" bson:"total_order_price"`
		CreatedAt       *time.Time       `json:"CreatedAt,omitempty" bson:"created_at"`
	}
	State     uint
	Side      uint
	Execution uint
)

func (o Order) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", fmt.Sprintf("%d", o.Id)).
		Str("orderId", fmt.Sprintf("%d", o.OrderId)).
		Str("execution", o.Execution.String()).
		Str("side", o.Side.String()).
		Str("srcCurrency", o.SrcCurrency.String()).
		Str("dstCurrency", o.DstCurrency.String()).
		Str("price", o.Price.String()).
		Str("amount", o.Amount.String()).
		Str("totalPrice", o.TotalPrice.String()).
		Str("totalOrderPrice", o.TotalOrderPrice.String()).
		Str("stopPrice", o.StopPrice.String()).
		Str("status", o.Status.String()).
		Str("fee", o.Fee.String()).
		Str("feeCurrency", o.FeeCurrency.String()).
		Str("user", o.Account).
		Str("createdAt", o.CreatedAt.String()).
		Str("unmatchedAmount", o.UnmatchedAmount.String()).
		Str("matchedAmount", o.MatchedAmount.String()).
		Bool("partial", o.Partial)
}

const (
	MarketExecution Execution = iota
	LimitExecution
	StopMarketExecution
	StopLimitExecution
)

var executionMap = map[string]Execution{
	"market":      MarketExecution,
	"limit":       LimitExecution,
	"stop_market": StopMarketExecution,
	"stop_limit":  StopLimitExecution,
}

func (o Execution) String() string {
	return string([]string{"market", "limit", "stop_market", "stop_limit"}[o])
}

func ExecutionFromString(s string) (Execution, error) {
	sym, ok := executionMap[s]
	if !ok {
		return 0, fmt.Errorf("invalid execution %s", s)
	}
	return sym, nil
}

const (
	Buy Side = iota
	Sell
)

var sideMap = map[string]Side{
	"buy":  Buy,
	"sell": Sell,
}

func (o Side) String() string {
	return []string{"buy", "sell"}[o]
}

func (o Side) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", o.String())), nil
}

func SideFromString(s string) (Side, error) {
	side, ok := sideMap[s]
	if !ok {
		return 0, fmt.Errorf("invalid side %s", s)
	}
	return side, nil
}

const (
	Open State = iota
	Filled
	PartiallyFilled
	Canceled
	Draft
)

var stateMap = map[string]State{
	"open":            Open,
	"filled":          Filled,
	"partiallyFilled": PartiallyFilled,
	"canceled":        Canceled,
	"draft":           Draft,
}

func StateFromString(s string) (State, error) {
	state, ok := stateMap[s]
	if !ok {
		return 0, fmt.Errorf("invalid state %s", s)
	}
	return state, nil
}

func (o State) String() string {
	return []string{"open", "filled", "partially_filled", "canceled", "draft"}[o]
}
