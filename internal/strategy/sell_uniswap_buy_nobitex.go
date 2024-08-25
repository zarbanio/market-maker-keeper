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
	"github.com/zarbanio/market-maker-keeper/internal/uniswapv3"
	"github.com/zarbanio/market-maker-keeper/pkg/logger"
	"github.com/zarbanio/market-maker-keeper/store"
)

type SellDaiUniswapBuyTetherNobitex struct {
	s             store.IStore
	nobitex       domain.Exchange
	uniswapQuoter *uniswapv3.Quoter
	dexTrader     *dextrader.Wrapper
	tokens        map[symbol.Symbol]domain.Token
	uniswapFee    domain.UniswapFee

	marketsdata MarketsData
	config      Config
}

func NewSellDaiUniswapBuyTetherNobitex(
	s store.IStore,
	exchange domain.Exchange,
	dexTrader *dextrader.Wrapper,
	uniswapQuoter *uniswapv3.Quoter,
	tokens map[symbol.Symbol]domain.Token,
	config Config,
) ArbitrageStrategy {
	marketsdata := make(map[Market]MarketData)
	marketsdata[UniswapV3] = NewMarketData()
	marketsdata[Nobitex] = NewMarketData()

	return &SellDaiUniswapBuyTetherNobitex{
		s:             s,
		nobitex:       exchange,
		uniswapQuoter: uniswapQuoter,
		dexTrader:     dexTrader,
		config:        config,
		tokens:        tokens,
		uniswapFee:    domain.UniswapFeeFee01,
		marketsdata:   marketsdata,
	}
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
	return s.marketsdata, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) getNobitexBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.nobitex.Balances()
	if err != nil {
		return nil, fmt.Errorf("failed to get nobitex balances. %w", err)
	}
	for _, balance := range balances {
		s.marketsdata[Nobitex].Balances[balance.Symbol] = balance.Balance
	}
	return s.marketsdata[Nobitex].Balances, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) getDexTraderBalances() (map[symbol.Symbol]decimal.Decimal, error) {
	balances, err := s.dexTrader.GetTokenBalances(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get dex trader balances. %w", err)
	}
	for _, balance := range balances {
		s.marketsdata[UniswapV3].Balances[balance.Symbol] = balance.Balance
	}
	return s.marketsdata[UniswapV3].Balances, nil
}

func (s *SellDaiUniswapBuyTetherNobitex) getNobitexPrices() (map[symbol.Symbol]decimal.Decimal, error) {
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

func (s *SellDaiUniswapBuyTetherNobitex) Evaluate(ctx context.Context) (*ArbitrageOpportunity, error) {
	if s.marketsdata == nil {
		return nil, fmt.Errorf("markets data is nil")
	}

	startQty := s.config.StartQty
	endQty := decimal.Min(s.config.EndQty)
	stepQty := s.config.StepQty

	var bestArbirageOpportunity *ArbitrageOpportunity

	for qty := startQty; qty.LessThanOrEqual(endQty); qty = qty.Add(stepQty) {
		uniswapV3OrderCandidate, err := s.findUniswapOrderCandidate(ctx, qty)
		if err != nil {
			if errors.Is(err, ErrorInsufficentBalance{}) {
				logger.Logger.Warn().Err(err).Msg("not enough balance to find the best arbitrage opportunity.")
				break
			}
			return nil, fmt.Errorf("failed to find uniswap order candidate. %w", err)
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
			return nil, fmt.Errorf("failed to find nobitex order candidate. %w", err)
		}

		arbitrageOpportunity := &ArbitrageOpportunity{
			UniV3OrderCandidate:   *uniswapV3OrderCandidate,
			NobitexOrderCandidate: *nobitexOrderCandidate,
		}

		logger.Logger.Debug().Str("strategy", s.Name()).Object("ArbitrageOpportunity", arbitrageOpportunity).Msg("arbitrage opportunity")

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

func (s *SellDaiUniswapBuyTetherNobitex) findUniswapOrderCandidate(ctx context.Context, qty decimal.Decimal) (*OrderCandidate, error) {
	// amountIn is in dai
	// amountOut is in zar
	// we are selling dai
	tokenIn := s.tokens[symbol.DAI]
	tokenOut := s.tokens[symbol.ZAR]
	etherPrice := s.marketsdata[Nobitex].Prices[symbol.ETH]
	usdtPrice := s.marketsdata[Nobitex].Prices[symbol.USDT]
	daiBalance := s.marketsdata[UniswapV3].Balances[symbol.DAI]

	out, err := s.uniswapQuoter.GetSwapOutputWithExactInput(ctx, tokenIn, tokenOut, s.uniswapFee, qty)
	if err != nil {
		return nil, fmt.Errorf("failed to get swap output with exact input: %w", err)
	}

	if !out.IsPositive() {
		return nil, errors.New("not valid quantity for this liquidity")
	}

	minOut := out.Mul(decimal.NewFromInt(1).Sub(s.config.Slippage)).Round(int32(tokenOut.Decimals()))

	dexTraderGasFee, err := s.dexTrader.EstimateDexTradeGasFee(tokenIn, tokenOut, s.uniswapFee.BigInt(), qty, minOut)
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

func (s *SellDaiUniswapBuyTetherNobitex) findNobitexOrderCandidate(ctx context.Context, qty decimal.Decimal) (*OrderCandidate, error) {
	// qty is in tether
	// makerFee := s.nobitex.Fees(trade.Maker)
	// minOrder := s.nobitex.MinimumOrderToman()
	usdtPrice := s.marketsdata[Nobitex].Prices[symbol.USDT]
	tomanBalance := s.marketsdata[Nobitex].Balances[symbol.RLS].Div(decimal.NewFromInt(10)) // convert from rial to toman

	orderBook, err := s.nobitex.OrderBook(symbol.USDT, symbol.IRT)
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
	s.marketsdata = marketsdata
}
