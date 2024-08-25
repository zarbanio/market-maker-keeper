package domain

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type UniswapFee int64

const (
	UniswapFee0005   UniswapFee = 500    // 0.05%
	UniswapFeeFee003 UniswapFee = 3_000  // 0.3%
	UniswapFeeFee01  UniswapFee = 10_000 // 1%
)

var (
	feesMap = map[float64]UniswapFee{
		0.0005: UniswapFee0005,
		0.003:  UniswapFeeFee003,
		0.01:   UniswapFeeFee01,
	}
)

func ParseUniswapFee(f float64) UniswapFee {
	return feesMap[f]
}

func (f UniswapFee) Decimal() decimal.Decimal {
	return decimal.New(int64(f), 0)
}

func (f UniswapFee) BigInt() *big.Int {
	return big.NewInt(int64(f))
}

func (f UniswapFee) Percent() decimal.Decimal {
	return decimal.NewFromFloat(float64(int64(f) / 10_00_000))
}
