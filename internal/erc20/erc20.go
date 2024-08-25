package erc20

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type Token struct {
	addr     common.Address
	symbol   symbol.Symbol
	decimals int64
}

func NewToken(addr common.Address, sym symbol.Symbol, decimals int64) Token {
	return Token{
		addr:     addr,
		symbol:   sym,
		decimals: decimals,
	}
}

func (t Token) Address() common.Address {
	return t.addr
}

func (t Token) Symbol() symbol.Symbol {
	return t.symbol
}

func (t Token) Decimals() int64 {
	return t.decimals
}
