package symbol

import (
	"fmt"
	"strings"
)

type Symbol string

const (
	DAI         Symbol = "DAI"
	ZAR         Symbol = "ZAR"
	BTC         Symbol = "BTC"
	ETH         Symbol = "ETH"
	USDT        Symbol = "USDT"
	USDC        Symbol = "USDC"
	RLS         Symbol = "RLS"
	IRT         Symbol = "IRT"
	TMN         Symbol = "TMN"
	PersianRial Symbol = "﷼"
	Tether      Symbol = "Tether"
)

type Nobitexable interface {
	Nobitexable() string
}

var (
	symbolsMap = map[string]Symbol{
		"DAI":  DAI,
		"ZAR":  ZAR,
		"BTC":  BTC,
		"ETH":  ETH,
		"USDT": USDT,
		"USDC": USDC,
		"RLS":  RLS,
		"IRT":  IRT,
		"TMN":  TMN,
		"﷼":    PersianRial,
	}
)

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) LowerCaseString() string {
	return strings.ToLower(s.String())
}

func FromString(s string) (Symbol, error) {
	if s == "" {
		return "", fmt.Errorf("symbol is empty")
	}
	if s == "Tether" {
		return Tether, nil
	}
	sym, ok := symbolsMap[strings.ToUpper(s)]
	if !ok {
		return "", fmt.Errorf("invalid symbol %s", s)
	}
	return sym, nil
}
