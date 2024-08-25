package domain

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type Balance struct {
	Symbol  symbol.Symbol
	Balance decimal.Decimal
}

func CommaSeparate(numStr string) string {
	// Split the number into integer and decimal parts
	parts := strings.Split(numStr, ".")

	// Convert the integer part to int64
	intPart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return numStr
	}

	// Format the integer part with commas
	intPartStr := strconv.FormatInt(intPart, 10)
	for i := len(intPartStr) - 3; i > 0; i -= 3 {
		intPartStr = intPartStr[:i] + "," + intPartStr[i:]
	}

	// If there is a decimal part, append it
	if len(parts) > 1 {
		intPartStr += "." + parts[1]
	}

	return intPartStr
}
