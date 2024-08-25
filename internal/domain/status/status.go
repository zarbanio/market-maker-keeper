package status

import (
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
)

type (
	Status struct {
		UnmatchedAmount decimal.Decimal
		MatchedAmount   decimal.Decimal
		Partial         bool
		State           order.State
	}
)

func NewStatus(UnmatchedAmount decimal.Decimal, MatchedAmount decimal.Decimal, Partial bool, State order.State) *Status {
	return &Status{
		UnmatchedAmount: UnmatchedAmount,
		MatchedAmount:   MatchedAmount,
		Partial:         Partial,
		State:           State,
	}
}
