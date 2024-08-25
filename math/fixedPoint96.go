package math

import "math/big"

var FixedPoint96 = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
