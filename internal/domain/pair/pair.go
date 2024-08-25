package pair

import (
	"fmt"

	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type (
	Pair struct {
		Id         int64         `json:"id"`
		BaseAsset  symbol.Symbol `json:"baseAsset"`
		QuoteAsset symbol.Symbol `json:"quoteAsset"`
	}
)

func (p Pair) Symbol() string {
	return fmt.Sprintf("%s%s", p.BaseAsset.String(), p.QuoteAsset.String())
}
