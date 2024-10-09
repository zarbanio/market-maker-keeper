package strategy

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/dextrader"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
	"github.com/zarbanio/market-maker-keeper/internal/uniswapv3"
	"github.com/zarbanio/market-maker-keeper/store"
)

var (
	ErrorInvalidPrice  = errors.New("invalid price")
	ErrorInvalidAmount = errors.New("invalid amount")
)

type ErrorInsufficentBalance struct {
	Symbol   symbol.Symbol
	Required decimal.Decimal
	Got      decimal.Decimal
}

func (e ErrorInsufficentBalance) Error() string {
	return fmt.Sprintf("insufficient %s balance. required: %s, got: %s", e.Symbol, domain.CommaSeparate(e.Required.String()), domain.CommaSeparate(e.Got.String()))
}

type BuyDaiUniswapSellTetherNobitex struct {
	Store      store.IStore
	Nobitex    domain.Exchange
	DexQuoter  *uniswapv3.Quoter
	DexTrader  *dextrader.Wrapper
	Tokens     map[symbol.Symbol]domain.Token
	UniswapFee domain.UniswapFee
	Logger     zerolog.Logger

	Marketsdata map[Market]MarketData
	Config      Config
}

func (s *BuyDaiUniswapSellTetherNobitex) Name() string {
	return "BuyDaiUniswapSellTetherNobitex"
}

func (s *BuyDaiUniswapSellTetherNobitex) Setup() (MarketsData, error) {
	_, err := s.getNobitexBalances()
	if err != nil {
		return nil, err
	}
	_, err = s.getDexTraderBalances()
	if err != nil {
		return nil, err
	}
	_, err = s.getNobitexPrices()
	if err != nil {
		return nil, err
	}
	return s.Marketsdata, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getNobitexBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.Nobitex.Balances()
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex balances. %w", err)
	}
	for _, balance := range balances {
		s.Marketsdata[Nobitex].Balances[balance.Symbol] = balance.Balance
	}
	return s.Marketsdata[Nobitex].Balances, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getDexTraderBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.DexTrader.GetTokenBalances(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get dex trader balances. %w", err)
	}
	for _, balance := range balances {
		if balance.Balance.IsZero() {
			continue
		}
		s.Marketsdata[UniswapV3].Balances[balance.Symbol] = balance.Balance
	}
	return s.Marketsdata[UniswapV3].Balances, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getNobitexPrices() (map[symbol.Symbol]decimal.Decimal, error) {
	ethPrice, err := s.Nobitex.ExchangeRate(symbol.ETH, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get eth price. %w", err)
	}
	tetherPrice, err := s.Nobitex.ExchangeRate(symbol.USDT, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get tether price. %w", err)
	}
	s.Marketsdata[Nobitex].Prices[symbol.ETH] = ethPrice.Div(decimal.NewFromInt(10))     // convert from rial to toman
	s.Marketsdata[Nobitex].Prices[symbol.USDT] = tetherPrice.Div(decimal.NewFromInt(10)) // convert from rial to toman
	return s.Marketsdata[Nobitex].Prices, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) Evaluate(ctx context.Context) (*ArbitrageOpportunity, error) {
	startQty := s.Config.StartQty // in dai
	endQty := decimal.Min(
		s.Marketsdata[Nobitex].Balances[symbol.USDT],
		s.Marketsdata[UniswapV3].Balances[symbol.DAI],
	)
	stepQty := s.Config.StepQty // in dai

	var bestArbirageOpportunity *ArbitrageOpportunity

	for qty := startQty; qty.LessThanOrEqual(endQty); qty = qty.Add(stepQty) {
		uniswapV3OrderCandidate, err := s.findUniswapOrderCandidate(ctx, qty)
		if err != nil {
			if errors.Is(err, ErrorInsufficentBalance{}) {
				break
			}
			return nil, err
		}

		nobitexOrderCandidate, err := s.findNobitexOrderCandidate(ctx, qty)
		if err != nil {
			if errors.Is(err, ErrorInsufficentBalance{}) {
				break
			}
			if errors.Is(err, ErrorInvalidAmount) {
				s.Logger.Warn().Err(err).Msg("invalid amount to find the best arbitrage opportunity.")
				continue
			}
			return nil, err
		}

		arbitrageOpportunity := &ArbitrageOpportunity{
			UniV3OrderCandidate:   *uniswapV3OrderCandidate,
			NobitexOrderCandidate: *nobitexOrderCandidate,
		}
		s.Logger.Debug().Object("ArbitrageOpportunity", arbitrageOpportunity).Msg("arbitrage opportunity")

		if arbitrageOpportunity.EstimatedProfit().GreaterThan(bestArbirageOpportunity.EstimatedProfit()) {
			bestArbirageOpportunity = arbitrageOpportunity
		}
	}

	if bestArbirageOpportunity == nil {
		return nil, nil
	}

	if !bestArbirageOpportunity.IsProfitable() {
		return nil, nil
	}

	return bestArbirageOpportunity, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) findUniswapOrderCandidate(ctx context.Context, qty decimal.Decimal) (*OrderCandidate, error) {
	// amountIn is in zar
	// amountOut is in dai
	// we are buying dai with zar
	tokenIn := s.Tokens[symbol.ZAR]
	tokenOut := s.Tokens[symbol.DAI]
	tetherPrice := s.Marketsdata[Nobitex].Prices[symbol.USDT]
	etherPrice := s.Marketsdata[Nobitex].Prices[symbol.ETH]
	zarBalance := s.Marketsdata[UniswapV3].Balances[symbol.ZAR]

	in, err := s.DexQuoter.GetSwapInputWithExactOutput(ctx, tokenIn, tokenOut, s.UniswapFee, qty)
	if err != nil {
		return nil, fmt.Errorf("failed to get swap output with exact output: %w", err)
	}

	dexTraderGasFee, err := s.DexTrader.EstimateDexTradeGasFee(tokenIn, tokenOut, s.UniswapFee.BigInt(), in, qty)
	if err != nil {
		// If the dexTrader's assets are less than the specified quantity,
		// the EstimateDexTradeGasFee function will return an "execution reverted: STF" error.
		// This error occurs when the execution is reverted by a require assertion
		// in the TransferHelper.safeTransferFrom function, indicating that the transfer failed.
		if err.Error() == "execution reverted: STF" {
			return nil, ErrorInsufficentBalance{
				Symbol:   symbol.ZAR,
				Required: in,
				Got:      zarBalance,
			}
		}
		return nil, fmt.Errorf("failed to estimate dex trade gas fee. %w", err)
	}

	txFeeValue := dexTraderGasFee.Mul(etherPrice)     // convert from ether to toman
	outgoingValue := in.Mul(decimal.NewFromInt(1000)) // convert from zar to toman

	cost := Cost{
		OutgoingValue: outgoingValue,
		NetworkFee:    txFeeValue,
	}

	revenue := qty.Mul(tetherPrice) // convert from dai to toman, we assume that the tether price is the same as the dai price

	// buying dai with zar
	// in is in zar
	// qty is in dai
	oc := &OrderCandidate{
		Side:   order.Buy,
		Pair:   DAI_ZAR,
		In:     in,
		MinOut: qty,
		Market: UniswapV3,

		EstimatedCost:    cost,
		EstimatedRevenue: revenue,
	}

	return oc, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) findNobitexOrderCandidate(ctx context.Context, qty decimal.Decimal) (*OrderCandidate, error) {
	// amountIn is in tether
	// amountOut is in toman
	// we are selling tether for toman
	takerFee := s.Nobitex.Fees(trade.Taker)
	minOrder := s.Nobitex.MinimumOrderToman()
	tetherPrice := s.Marketsdata[Nobitex].Prices[symbol.USDT]
	tetherBalance := s.Marketsdata[Nobitex].Balances[symbol.USDT]

	orderBook, err := s.Nobitex.OrderBook(symbol.USDT, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex order book: %w", err)
	}

	in := qty.Mul(tetherPrice) // convert from tether to toman
	if in.LessThanOrEqual(minOrder) {
		return nil, fmt.Errorf("nobitex order amount is less than minimum order amount. %w", ErrorInvalidAmount)
	}
	if tetherBalance.LessThanOrEqual(qty) {
		return nil, ErrorInsufficentBalance{
			Symbol:   symbol.USDT,
			Required: qty,
			Got:      tetherBalance,
		}
	}

	avgPrice := orderbook.GetOrderBookPrice(orderBook, order.Sell, qty).Div(decimal.NewFromInt(10)) // convert from rial to toman
	if !avgPrice.IsPositive() {
		return nil, fmt.Errorf("invalid price. %w", ErrorInvalidPrice)
	}

	out := avgPrice.Mul(qty)
	out = out.Div(decimal.NewFromInt(1).Sub(takerFee))

	oc := &OrderCandidate{
		Side:   order.Sell,
		Pair:   USDT_TMN,
		In:     qty,
		MinOut: out,
		Market: Nobitex,

		EstimatedCost: Cost{
			OutgoingValue: in,
			NetworkFee:    decimal.Zero,
		},
		EstimatedRevenue: out,
	}

	return oc, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) Teardown() {
	marketsdata := make(map[Market]MarketData)
	marketsdata[UniswapV3] = NewMarketData()
	marketsdata[Nobitex] = NewMarketData()
	s.Marketsdata = marketsdata
}
