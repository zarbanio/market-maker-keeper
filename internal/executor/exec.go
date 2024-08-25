package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zarbanio/market-maker-keeper/internal/chain"
	"github.com/zarbanio/market-maker-keeper/internal/dextrader"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"
	"github.com/zarbanio/market-maker-keeper/internal/strategy"
	"github.com/zarbanio/market-maker-keeper/pkg/logger"
	"github.com/zarbanio/market-maker-keeper/store"
)

type Executor struct {
	s                    store.IStore
	strategies           []strategy.ArbitrageStrategy
	nobitex              domain.Exchange
	dexTrader            dextrader.Wrapper
	indxer               chain.Indexer
	cycleId              int64
	pairId               int64
	orderId              int64
	transactionId        int64
	nobitexRetryTimeOut  time.Duration
	nobitexSleepDuration time.Duration
	uniswapFee           domain.UniswapFee
}

func NewExecutor(
	s store.IStore,
	pairId int64,
	strategies []strategy.ArbitrageStrategy,
	nobitex domain.Exchange,
	dexTrader dextrader.Wrapper,
	indxer chain.Indexer,
	nobitexRetryTimeOut time.Duration,
	nobitexSleepDuration time.Duration,

) *Executor {
	return &Executor{
		s:                    s,
		pairId:               pairId,
		strategies:           strategies,
		nobitex:              nobitex,
		dexTrader:            dexTrader,
		indxer:               indxer,
		nobitexRetryTimeOut:  nobitexRetryTimeOut,
		nobitexSleepDuration: nobitexSleepDuration,
		uniswapFee:           domain.UniswapFeeFee01,
	}
}

func (e *Executor) RunAll() {
	for _, strategy := range e.strategies {
		err := e.Run(strategy)
		if err != nil {
			logger.Logger.Error().Err(err).Str("strategy", strategy.Name()).Msg("failed to run strategy")
		}
	}
}

func (e *Executor) Run(strategy strategy.ArbitrageStrategy) error {
	marketdata, err := strategy.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup strategy %s. %w", strategy.Name(), err)
	}
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategy.Name()).Object("marketdata", marketdata).Msg("strategy setup completed.")

	opportunity, err := strategy.Evaluate(context.Background())
	if err != nil {
		return fmt.Errorf("failed to evaluate strategy %s. %w", strategy.Name(), err)
	}
	if opportunity == nil {
		logger.Logger.Info().
			Int64("cycleId", e.cycleId).
			Str("strategy", strategy.Name()).
			Msg("no profitable opportunity found.")
		return nil
	}

	logger.Logger.Info().
		Int64("cycleId", e.cycleId).
		Str("strategy", strategy.Name()).
		Object("bestArbirageOpportunity", opportunity).
		Msg("a profitable opportunity found.")

	err = e.Execute(strategy.Name(), opportunity.NobitexOrderCandidate)
	if err != nil {
		return fmt.Errorf("failed to execute strategy %s. %w", strategy.Name(), err)
	}

	err = e.Execute(strategy.Name(), opportunity.UniV3OrderCandidate)
	if err != nil {
		return fmt.Errorf("failed to execute strategy %s. %w", strategy.Name(), err)
	}

	_, err = e.s.CreateNewTrade(context.Background(), e.pairId, e.orderId, e.transactionId)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&opportunity)
	if err != nil {
		return err
	}

	err = e.s.CreateArbitrageOpporchunity(context.Background(), data)
	if err != nil {
		return err
	}

	strategy.Teardown()

	return nil
}

func (e *Executor) Execute(strategyName string, orderCandidate strategy.OrderCandidate) error {
	switch orderCandidate.Market {
	case strategy.UniswapV3:
		return e.ExecuteUniswap(strategyName, orderCandidate)
	case strategy.Nobitex:
		return e.ExecuteNobitex(strategyName, orderCandidate)
	}
	return fmt.Errorf("invalid market: %s", orderCandidate.Market)
}

func (e *Executor) SetCycleId(cycleId int64) {
	e.cycleId = cycleId
}

func (e *Executor) ExecuteUniswap(strategyName string, orderCandidate strategy.OrderCandidate) error {
	var (
		txID   int64
		txHash common.Hash
	)
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Object("orderCandidate", orderCandidate).Msg("executing uniswap order candidate")

	tx, err := e.dexTrader.Trade(orderCandidate.Source(), orderCandidate.Destination(), e.uniswapFee, orderCandidate.In, orderCandidate.MinOut)
	if err != nil {
		return err
	}
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("transaction sent to uniswap.")
	txID, err = e.s.CreateTransaction(context.Background(), tx, e.dexTrader.GetExecutorAddress())
	if err != nil {
		return err
	}
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Int64("transaction_id", txID).Msg("transaction created in database.")

	txHash = tx.Hash()
	e.transactionId = txID

	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("waiting for transaction receipt.")
	rec, header, err := e.indxer.WaitForReceipt(context.Background(), txHash)
	if err != nil {
		return err
	}
	txUpdate := transaction.UpdatedFields{
		GasUsed:     rec.GasUsed,
		BlockNumber: rec.BlockNumber.Int64(),
		Timestamp:   time.Unix(int64(header.Time), 0),
		Status:      transaction.CastFromReceiptStatus(rec.Status),
	}
	err = e.s.UpdateTransaction(context.Background(), txID, txUpdate)
	if err != nil {
		return err
	}
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("transaction receipt received.")
	return nil
}

func (e *Executor) ExecuteNobitex(strategyName string, orderCandidate strategy.OrderCandidate) error {
	var (
		nobitexOrderId int64
		orderId        int64
	)

	o := order.Order{
		Side:        orderCandidate.Side,
		Execution:   order.MarketExecution,
		SrcCurrency: orderCandidate.Source(),
		DstCurrency: orderCandidate.Destination(),
		FeeCurrency: symbol.TMN,
		Amount:      orderCandidate.Amount(),
		Status:      order.Draft,
	}

	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Object("order", o).Msg("executing nobitex order candidate")

	id, t, err := e.nobitex.PlaceOrder(o)
	if err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	if id == 0 {
		return fmt.Errorf("invalid order id: %d", id)
	}

	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Int64("nobitexOrderId", id).Msg("order placed in nobitex.")

	nobitexOrderId = id
	o.CreatedAt = t
	o.OrderId = nobitexOrderId
	o.Status = order.Open
	orderID, err := e.s.CreateNewOrder(context.Background(), o)
	orderId = orderID
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Int64("orderId", orderId).Msg("order created in database.")

	e.orderId = orderID

	logger.Logger.Info().Int64("cycleId", e.cycleId).Str("strategy", strategyName).Int64("nobitexOrderId", nobitexOrderId).Msg("waiting for nobitex order status.")
	ctx, cancel := context.WithTimeout(context.Background(), e.nobitexRetryTimeOut)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.Logger.Error().Int64("nobitexOrderId", nobitexOrderId).Msg("nobitex order status timeout")
			return ctx.Err()
		default:
			orderUpdate, err := e.nobitex.OrderStatus(context.Background(), nobitexOrderId)
			if err != nil {
				return err
			}
			if orderUpdate.Status != order.Filled {
				continue
			}
			logger.Logger.Info().Int64("nobitexOrderId", nobitexOrderId).Str("status", orderUpdate.Status.String()).Msg("nobitex order status received.")
			update := order.UpdatedFields{
				Status:          &orderUpdate.Status,
				Price:           &orderUpdate.Price,
				TotalOrderPrice: &orderUpdate.TotalOrderPrice,
				TotalPrice:      &orderUpdate.TotalPrice,
				Fee:             &orderUpdate.Fee,
				MatchedAmount:   &orderUpdate.MatchedAmount,
				UnmatchedAmount: &orderUpdate.UnmatchedAmount,
				CreatedAt:       &orderUpdate.CreatedAt,
			}
			err = e.s.UpdateOrder(context.Background(), orderId, update)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
