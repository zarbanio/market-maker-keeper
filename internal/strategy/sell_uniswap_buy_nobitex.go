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
	"github.com/zarbanio/market-maker-keeper/internal/uniswapv3"
	"github.com/zarbanio/market-maker-keeper/store"
)

type SellDaiUniswapBuyTetherNobitex struct {
	Store         store.IStore
	Nobitex       domain.Exchange
	UniswapQuoter *uniswapv3.Quoter
	DexTrader     *dextrader.Wrapper
	Tokens        map[symbol.Symbol]domain.Token
	UniswapFee    domain.UniswapFee
	Logger        zerolog.Logger

	Marketsdata MarketsData
	Config      Config
}

func (s *SellDaiUniswapBuyTetherNobitex) Name() string {
	return "SellDaiUniswapBuyTetherNobitex"
}

func (s *SellDaiUniswapBuyTetherNobitex) Setup() (MarketsData, error) {
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

func (s *SellDaiUniswapBuyTetherNobitex) getNobitexBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.Nobitex.Balances()
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex balances. %w", err)
	}
	for _, balance := range balances {
		if balance.Balance.IsZero() {
			continue
		}
		s.Marketsdata[Nobitex].Balances[balance.Symbol] = balance.Balance
	}
	return s.Marketsdata[Nobitex].Balances, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) getDexTraderBalances() (map[symbol.Symbol]decimal.Decimal, error) {
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

func (s *SellDaiUniswapBuyTetherNobitex) getNobitexPrices() (map[symbol.Symbol]decimal.Decimal, error) {
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

func (s *SellDaiUniswapBuyTetherNobitex) Evaluate(ctx context.Context) (*ArbitrageOpportunity, error) {
	if s.Marketsdata == nil {
		return nil, fmt.Errorf("markets data is nil")
	}

	startQty := s.Config.StartQty
	endQty := s.Marketsdata[UniswapV3].Balances[symbol.DAI]
	stepQty := s.Config.StepQty

	var bestArbirageOpportunity *ArbitrageOpportunity

	s.Logger.Debug().
		Str("startQty", domain.CommaSeparate(startQty.String())).
		Str("endQty", domain.CommaSeparate(endQty.String())).
		Str("stepQty", domain.CommaSeparate(stepQty.String())).
		Str("strategy", s.Name()).
		Msg("evaluating arbitrage opportunity")

	for qty := startQty; qty.LessThanOrEqual(endQty); qty = qty.Add(stepQty) {
		s.Logger.Debug().
			Str("qty", domain.CommaSeparate(qty.String())).
			Msg("evaluating arbitrage opportunity")

		uniswapV3OrderCandidate, err := s.findUniswapOrderCandidate(ctx, qty)
		if err != nil {
			if errors.As(err, &ErrorInsufficentBalance{}) {
				break
			}
			return nil, fmt.Errorf("failed to find uniswap order candidate. %w", err)
		}

		s.Logger.Debug().
			Object("uniswapV3OrderCandidate", uniswapV3OrderCandidate).
			Msg("uniswap v3 order candidate found")

		nobitexOrderCandidate, err := s.findNobitexOrderCandidate(qty)
		if err != nil {
			if errors.As(err, &ErrorInsufficentBalance{}) {
				break
			}
			if errors.Is(err, ErrorInvalidAmount) {
				s.Logger.Warn().Err(err).Msg("invalid amount to find the best arbitrage opportunity.")
				continue
			}
			return nil, fmt.Errorf("failed to find nobitex order candidate. %w", err)
		}

		s.Logger.Debug().
			Object("nobitexOrderCandidate", nobitexOrderCandidate).
			Msg("nobitex order candidate found")

		arbitrageOpportunity := &ArbitrageOpportunity{
			UniV3OrderCandidate:   *uniswapV3OrderCandidate,
			NobitexOrderCandidate: *nobitexOrderCandidate,
		}

		s.Logger.Debug().Str("strategy", s.Name()).Object("ArbitrageOpportunity", arbitrageOpportunity).Msg("arbitrage opportunity")

		if arbitrageOpportunity.EstimatedProfit().GreaterThan(bestArbirageOpportunity.EstimatedProfit()) {
			bestArbirageOpportunity = arbitrageOpportunity

			s.Logger.Debug().
				Object("arbitrageOpportunity", bestArbirageOpportunity).
				Msg("best arbitrage opportunity so far")
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

func (s *SellDaiUniswapBuyTetherNobitex) findUniswapOrderCandidate(ctx context.Context, qty decimal.Decimal) (*OrderCandidate, error) {
	// amountIn is in dai
	// amountOut is in zar
	// we are selling dai
	tokenIn := s.Tokens[symbol.DAI]
	tokenOut := s.Tokens[symbol.ZAR]
	etherPrice := s.Marketsdata[Nobitex].Prices[symbol.ETH]
	usdtPrice := s.Marketsdata[Nobitex].Prices[symbol.USDT]
	daiBalance := s.Marketsdata[UniswapV3].Balances[symbol.DAI]

	out, err := s.UniswapQuoter.GetSwapOutputWithExactInput(ctx, tokenIn, tokenOut, s.UniswapFee, qty)
	if err != nil {
		return nil, fmt.Errorf("failed to get swap output with exact input: %w", err)
	}

	if !out.IsPositive() {
		return nil, errors.New("not valid quantity for this liquidity")
	}

	minOut := out.Mul(decimal.NewFromInt(1).Sub(s.Config.Slippage)).Round(int32(tokenOut.Decimals()))

	dexTraderGasFee, err := s.DexTrader.EstimateDexTradeGasFee(tokenIn, tokenOut, s.UniswapFee.BigInt(), qty, minOut)
	if err != nil {
		// If the dexTrader's assets are less than the specified quantity,
		// the EstimateDexTradeGasFee function will return an "execution reverted: STF" error.
		// This error occurs when the execution is reverted by a require assertion
		// in the TransferHelper.safeTransferFrom function, indicating that the transfer failed.
		if err.Error() == "execution reverted: STF" {
			return nil, ErrorInsufficentBalance{
				Symbol:   symbol.DAI,
				Required: qty,
				Got:      daiBalance,
			}
		}
		return nil, fmt.Errorf("failed to estimate dex trade gas fee. %w", err)
	}

	txFeeValue := dexTraderGasFee.Mul(etherPrice) // convert from ether to toman
	outgoingValue := qty.Mul(usdtPrice)           // we multiply by usdt price because we are selling dai, we assume that the tether price is the same as the dai price

	cost := Cost{
		OutgoingValue: outgoingValue,
		NetworkFee:    txFeeValue,
	}

	revenue := out.Mul(decimal.NewFromInt(1000)) // we multiply by 1000 out is in zar and we want to convert it to toman

	// selling dai to zar
	// in is in dai
	// qty is in zar
	oc := &OrderCandidate{
		Side:             order.Sell,
		Pair:             DAI_ZAR,
		In:               qty,
		Out:              out,
		MinOut:           minOut,
		Market:           UniswapV3,
		EstimatedCost:    cost,
		EstimatedRevenue: revenue,
	}

	return oc, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) findNobitexOrderCandidate(qty decimal.Decimal) (*OrderCandidate, error) {
	// qty is in tether
	// makerFee := s.nobitex.Fees(trade.Maker)
	// minOrder := s.nobitex.MinimumOrderToman()
	usdtPrice := s.Marketsdata[Nobitex].Prices[symbol.USDT]
	tomanBalance := s.Marketsdata[Nobitex].Balances[symbol.RLS].Div(decimal.NewFromInt(10)) // convert from rial to toman

	orderBook, err := s.Nobitex.OrderBook(symbol.USDT, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex order book: %w", err)
	}

	avgPrice := orderbook.GetOrderBookPrice(orderBook, order.Buy, qty).Div(decimal.NewFromInt(10)) // convert from rial to toman
	if !avgPrice.IsPositive() {
		return nil, fmt.Errorf("invalid price. %w", ErrorInvalidPrice)
	}

	tomanRequired := avgPrice.Mul(qty)

	if tomanRequired.GreaterThan(tomanBalance) {
		return nil, ErrorInsufficentBalance{
			Symbol:   symbol.TMN,
			Required: tomanRequired,
			Got:      tomanBalance,
		}
	}

	revenue := qty.Mul(usdtPrice)

	oc := &OrderCandidate{
		Side:   order.Buy,
		Pair:   USDT_TMN,
		In:     tomanRequired,
		Out:    qty,
		MinOut: qty,
		Market: Nobitex,

		EstimatedCost: Cost{
			OutgoingValue: tomanRequired,
			NetworkFee:    decimal.Zero,
		},
		EstimatedRevenue: revenue,
	}

	return oc, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) Teardown() {
	marketsdata := make(map[Market]MarketData)
	marketsdata[UniswapV3] = NewMarketData()
	marketsdata[Nobitex] = NewMarketData()
	s.Marketsdata = marketsdata
}
