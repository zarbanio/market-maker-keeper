package math

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/cockroachdb/apd"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/shopspring/decimal"
)

var cache map[int64]*big.Int
var once sync.Once
var baseContext = apd.BaseContext

func Pow10(n int64) *big.Int {
	once.Do(func() {
		cache = make(map[int64]*big.Int)
	})
	if v, ok := cache[n]; ok {
		return v
	}
	cache[n] = math.BigPow(10, n)
	return cache[n]
}

func Normalize(n *big.Int, decimals int64) decimal.Decimal {
	return decimal.NewFromBigInt(n, 0).Div(decimal.NewFromBigInt(Pow10(decimals), 0))
}

func Denormalize(n decimal.Decimal, decimals int64) *big.Int {
	return n.Mul(decimal.New(10, int32(decimals-1))).BigInt()
}

func BigIntFromString(str string) *big.Int {
	ret, _ := new(big.Int).SetString(str, 10)
	return ret
}

func DecimalToBigInt(n decimal.Decimal, decimals int32) (*big.Int, error) {
	c := baseContext.WithPrecision(100)

	newN, condition, err := apd.NewFromString(n.String())
	if condition != 0 {
		return nil, fmt.Errorf(condition.String())
	}
	if err != nil {
		return nil, err
	}

	r := new(apd.Decimal)
	condition, err = c.Mul(r, newN, apd.New(1, decimals))
	if condition != 0 {
		return nil, fmt.Errorf(condition.String())
	}
	if err != nil {
		return nil, err
	}

	result, ok := new(big.Int).SetString(r.Text('f'), 10)
	if ok == false {
		return nil, fmt.Errorf("error in convert string to big.Int")
	}

	return result, nil
}

func BigIntToDecimal(n *big.Int, decimals int32) (decimal.Decimal, error) {
	c := baseContext.WithPrecision(100)

	newN, condition, err := apd.NewFromString(n.String())
	if condition != 0 {
		return decimal.Zero, fmt.Errorf(condition.String())
	}
	if err != nil {
		return decimal.Zero, err
	}

	r := new(apd.Decimal)
	condition, err = c.Quo(r, newN, apd.New(1, decimals))
	if condition != 0 {
		return decimal.Zero, fmt.Errorf(condition.String())
	}
	if err != nil {
		return decimal.Zero, err
	}

	result, err := decimal.NewFromString(r.String())
	if err != nil {
		return decimal.Zero, err
	}

	return result, nil
}
