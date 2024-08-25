package domain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type Token interface {
	Address() common.Address
	Symbol() symbol.Symbol
	Decimals() int64
}
