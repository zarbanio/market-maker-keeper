package strategy

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/dextrader"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
	"github.com/zarbanio/market-maker-keeper/internal/uniswapv3"
	"github.com/zarbanio/market-maker-keeper/pkg/logger"
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
	s          store.IStore
	nobitex    domain.Exchange
	dexQuoter  *uniswapv3.Quoter
	dexTrader  *dextrader.Wrapper
	tokens     map[symbol.Symbol]domain.Token
	uniswapFee domain.UniswapFee

	marketsdata map[Market]MarketData
	config      Config
}

func NewBuyDaiUniswapSellTetherNobitex(
	s store.IStore,
	exchange domain.Exchange,
	dexTrader *dextrader.Wrapper,
	dexQuoter *uniswapv3.Quoter,
	tokens map[symbol.Symbol]domain.Token,
	config Config,
) ArbitrageStrategy {
	marketsdata := make(map[Market]MarketData)
	marketsdata[UniswapV3] = NewMarketData()
	marketsdata[Nobitex] = NewMarketData()

	return &BuyDaiUniswapSellTetherNobitex{
		s:           s,
		uniswapFee:  domain.UniswapFeeFee01,
		nobitex:     exchange,
		dexQuoter:   dexQuoter,
		dexTrader:   dexTrader,
		config:      config,
		tokens:      tokens,
		marketsdata: marketsdata,
	}
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
	return s.marketsdata, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getNobitexBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.nobitex.Balances()
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex balances. %w", err)
	}
	for _, balance := range balances {
		s.marketsdata[Nobitex].Balances[balance.Symbol] = balance.Balance
	}
	return s.marketsdata[Nobitex].Balances, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getDexTraderBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.dexTrader.GetTokenBalances(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get dex trader balances. %w", err)
	}
	for _, balance := range balances {
		s.marketsdata[UniswapV3].Balances[balance.Symbol] = balance.Balance
	}
	return s.marketsdata[UniswapV3].Balances, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) getNobitexPrices() (map[symbol.Symbol]decimal.Decimal, error) {
	ethPrice, err := s.nobitex.ExchangeRate(symbol.ETH, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get eth price. %w", err)
	}
	tetherPrice, err := s.nobitex.ExchangeRate(symbol.USDT, symbol.IRT)
	if err != nil {
		return nil, fmt.Errorf("failed to get tether price. %w", err)
	}
	s.marketsdata[Nobitex].Prices[symbol.ETH] = ethPrice.Div(decimal.NewFromInt(10))     // convert from rial to toman
	s.marketsdata[Nobitex].Prices[symbol.USDT] = tetherPrice.Div(decimal.NewFromInt(10)) // convert from rial to toman
	return s.marketsdata[Nobitex].Prices, nil
}

func (s *BuyDaiUniswapSellTetherNobitex) Evaluate(ctx context.Context) (*ArbitrageOpportunity, error) {
	if s.marketsdata == nil {
		return nil, fmt.Errorf("markets data is nil")
	}

	startQty := s.config.StartQty                                                        // in dai
	endQty := decimal.Min(s.config.EndQty, s.marketsdata[Nobitex].Balances[symbol.USDT]) // in dai, we can't buy more than what we have in tether in nobitex
	stepQty := s.config.StepQty                                                          // in dai

	var bestArbirageOpportunity *ArbitrageOpportunity

	for qty := startQty; qty.LessThanOrEqual(endQty); qty = qty.Add(stepQty) {
		uniswapV3OrderCandidate, err := s.findUniswapOrderCandidate(ctx, qty)
		if err != nil {
			if errors.Is(err, ErrorInsufficentBalance{}) {
				logger.Logger.Warn().Err(err).Msg("not enough balance to find the best arbitrage opportunity.")
				break
			}
			return nil, err
		}

		nobitexOrderCandidate, err := s.findNobitexOrderCandidate(ctx, qty)
		if err != nil {
			if errors.Is(err, ErrorInsufficentBalance{}) {
				logger.Logger.Warn().Err(err).Msg("not enough balance to find the best arbitrage opportunity.")
				break
			}
			if errors.Is(err, ErrorInvalidAmount) {
				logger.Logger.Warn().Err(err).Msg("invalid amount to find the best arbitrage opportunity.")
				continue
			}
			return nil, err
		}

		arbitrageOpportunity := &ArbitrageOpportunity{
			UniV3OrderCandidate:   *uniswapV3OrderCandidate,
			NobitexOrderCandidate: *nobitexOrderCandidate,
		}
		logger.Logger.Debug().Object("ArbitrageOpportunity", arbitrageOpportunity).Msg("arbitrage opportunity")

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
	tokenIn := s.tokens[symbol.ZAR]
	tokenOut := s.tokens[symbol.DAI]
	tetherPrice := s.marketsdata[Nobitex].Prices[symbol.USDT]
	etherPrice := s.marketsdata[Nobitex].Prices[symbol.ETH]
	zarBalance := s.marketsdata[UniswapV3].Balances[symbol.ZAR]

	in, err := s.dexQuoter.GetSwapInputWithExactOutput(ctx, tokenIn, tokenOut, s.uniswapFee, qty)
	if err != nil {
		return nil, fmt.Errorf("failed to get swap output with exact output: %w", err)
	}

	dexTraderGasFee, err := s.dexTrader.EstimateDexTradeGasFee(tokenIn, tokenOut, s.uniswapFee.BigInt(), in, qty)
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
	takerFee := s.nobitex.Fees(trade.Taker)
	minOrder := s.nobitex.MinimumOrderToman()
	tetherPrice := s.marketsdata[Nobitex].Prices[symbol.USDT]
	tetherBalance := s.marketsdata[Nobitex].Balances[symbol.USDT]

	orderBook, err := s.nobitex.OrderBook(symbol.USDT, symbol.IRT)
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
	s.marketsdata = marketsdata
}
