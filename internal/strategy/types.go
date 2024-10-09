package strategy

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type Config struct {
	StartQty        decimal.Decimal
	StepQty         decimal.Decimal
	ProfitThreshold decimal.Decimal
	Slippage        decimal.Decimal
}

type Market string

const (
	UniswapV3 Market = "UniswapV3"
	Nobitex   Market = "Nobitex"
)

type Pair string

const (
	DAI_ZAR  Pair = "DAI/ZAR"
	USDT_TMN Pair = "USDT/TMN"
)

func (p Pair) Token0() symbol.Symbol {
	tokens := strings.Split(string(p), "/")
	return symbol.Symbol(tokens[0])
}

func (p Pair) Token1() symbol.Symbol {
	tokens := strings.Split(string(p), "/")
	return symbol.Symbol(tokens[1])
}

type OrderCandidate struct {
	Side   order.Side
	Pair   Pair
	In     decimal.Decimal
	Out    decimal.Decimal
	MinOut decimal.Decimal
	Market Market

	// EstimatedCost and EstimatedRevenue are used to calculate the profit/loss of the arbitrage opportunity.
	// Both are in the same currency, which is Tomans
	EstimatedCost    Cost
	EstimatedRevenue decimal.Decimal
}

type Cost struct {
	// NetworkFee is the fee that is paid to the network for the transaction
	NetworkFee decimal.Decimal

	// OutgoingValue is the value that is sent to the destination
	OutgoingValue decimal.Decimal
}

func (c Cost) Total() decimal.Decimal {
	return c.NetworkFee.Add(c.OutgoingValue)
}

func (c Cost) String() string {
	return "NetworkFee: " + c.NetworkFee.String() + "\n" + "OutgoingValue: " + c.OutgoingValue.String() + "\n" + "Total: " + c.Total().String() + "\n"
}

func (c Cost) MarshalJSON() ([]byte, error) {
	type Alias Cost
	return json.Marshal(&struct {
		Total string `json:"Total"`
		*Alias
	}{
		Total: c.Total().String(),
		Alias: (*Alias)(&c),
	})
}

type MarketData struct {
	Prices   map[symbol.Symbol]decimal.Decimal
	Balances map[symbol.Symbol]decimal.Decimal
}

type MarketsData map[Market]MarketData

func (m MarketsData) MarshalZerologObject(e *zerolog.Event) {
	for key, value := range m {
		e.Dict(string(key), zerolog.Dict().Fields(map[string]interface{}{
			"Prices":   value.Prices,
			"Balances": value.Balances,
		}))
	}
}

func NewMarketData() MarketData {
	return MarketData{
		Prices:   make(map[symbol.Symbol]decimal.Decimal),
		Balances: make(map[symbol.Symbol]decimal.Decimal),
	}
}

func (o *OrderCandidate) EstimatedProfit() decimal.Decimal {
	return o.EstimatedRevenue.Sub(o.EstimatedCost.Total())
}

type ArbitrageOpportunity struct {
	UniV3OrderCandidate   OrderCandidate
	NobitexOrderCandidate OrderCandidate
}

func (a ArbitrageOpportunity) String() string {
	return a.NobitexOrderCandidate.String() + "\n" + a.UniV3OrderCandidate.String() + "\n" + "EstimatedProfit: " + a.EstimatedProfit().String() + "\n"
}

func (a ArbitrageOpportunity) MarshalJSON() ([]byte, error) {
	type Alias ArbitrageOpportunity
	return json.Marshal(&struct {
		EstimatedProfit string `json:"EstimatedProfit"`
		*Alias
	}{
		EstimatedProfit: a.EstimatedProfit().String(),
		Alias:           (*Alias)(&a),
	})
}

func (a ArbitrageOpportunity) MarshalZerologObject(e *zerolog.Event) {
	e.Object("NobitexOrderCandidate", a.NobitexOrderCandidate)
	e.Object("UniV3OrderCandidate", a.UniV3OrderCandidate)
	e.Str("EstimatedProfit", a.EstimatedProfit().String())
}

func (oc OrderCandidate) MarshalZerologObject(e *zerolog.Event) {
	e.Str("Market", string(oc.Market))
	e.Str("Side", oc.Side.String())
	e.Str("Pair", string(oc.Pair))
	e.Str("In", domain.CommaSeparate(oc.In.String()))
	e.Str("MinOut", domain.CommaSeparate(oc.MinOut.String()))
	e.Object("EstimatedCost", oc.EstimatedCost)
	e.Str("EstimatedRevenue", domain.CommaSeparate(oc.EstimatedRevenue.String()))
}

func (c Cost) MarshalZerologObject(e *zerolog.Event) {
	e.Str("NetworkFee", domain.CommaSeparate(c.NetworkFee.String()))
	e.Str("OutgoingValue", domain.CommaSeparate(c.OutgoingValue.String()))
	e.Str("Total", domain.CommaSeparate(c.Total().String()))
}

func (oc OrderCandidate) String() string {
	return string(oc.Market) + "\n" +
		"Side: " + oc.Side.String() + "\n" +
		"Pair: " + string(oc.Pair) + "\n" +
		"In: " + oc.In.String() + "\n" +
		"MinOut: " + oc.MinOut.String() + "\n" +
		"EstimatedCost: " + "\n" + oc.EstimatedCost.String() + "\n" +
		"EstimatedRevenue: " + oc.EstimatedRevenue.String() + "\n"
}

func (oc OrderCandidate) Source() symbol.Symbol {
	if oc.Market == Nobitex {
		return oc.Pair.Token0()
	}
	if oc.Side == order.Buy {
		return oc.Pair.Token1()
	}
	return oc.Pair.Token0()
}

func (oc OrderCandidate) Amount() decimal.Decimal {
	if oc.Side == order.Buy {
		return oc.MinOut
	}
	return oc.In
}

func (oc OrderCandidate) Destination() symbol.Symbol {
	if oc.Market == Nobitex {
		return oc.Pair.Token1()
	}
	if oc.Side == order.Buy {
		return oc.Pair.Token0()
	}
	return oc.Pair.Token1()
}

func (a *ArbitrageOpportunity) EstimatedProfit() decimal.Decimal {
	if a == nil {
		return decimal.NewFromInt(math.MinInt64)
	}
	return a.NobitexOrderCandidate.EstimatedProfit().Add(a.UniV3OrderCandidate.EstimatedProfit())
}

func (a *ArbitrageOpportunity) IsProfitable() bool {
	return a.EstimatedProfit().GreaterThan(decimal.Zero)
}

type ArbitrageStrategy interface {
	Name() string
	Setup() (MarketsData, error)
	Evaluate(ctx context.Context) (*ArbitrageOpportunity, error)
	Teardown()
}
