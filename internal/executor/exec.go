package executor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/zarbanio/market-maker-keeper/internal/chain"
	"github.com/zarbanio/market-maker-keeper/internal/dextrader"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"
	"github.com/zarbanio/market-maker-keeper/internal/strategy"
	"github.com/zarbanio/market-maker-keeper/store"
)

type Executor struct {
	Store                store.IStore
	Strategies           []strategy.ArbitrageStrategy
	Nobitex              domain.Exchange
	DexTrader            dextrader.Wrapper
	Indxer               chain.Indexer
	CycleId              int64
	PairId               int64
	OrderId              int64
	TransactionId        int64
	NobitexRetryTimeOut  time.Duration
	NobitexSleepDuration time.Duration
	UniswapFee           domain.UniswapFee
	Logger               zerolog.Logger
}

func (e *Executor) RunAll() {
	for _, strategy := range e.Strategies {
		err := e.Run(strategy)
		if err != nil {
			e.Logger.Error().Err(err).Str("strategy", strategy.Name()).Msg("failed to run strategy")
		}
	}
}

func (e *Executor) Run(strg strategy.ArbitrageStrategy) error {
	marketdata, err := strg.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup strategy %s. %w", strg.Name(), err)
	}
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strg.Name()).Object("marketdata", marketdata).Msg("strategy setup completed.")

	opportunity, err := strg.Evaluate(context.Background())
	if err != nil {
		if !errors.As(err, &strategy.ErrorInsufficentBalance{}) {
			return fmt.Errorf("failed to evaluate strategy %s. %w", strg.Name(), err)
		}
	}
	if opportunity == nil {
		e.Logger.Info().
			Int64("cycleId", e.CycleId).
			Str("strategy", strg.Name()).
			Msg("no profitable opportunity found.")
		return nil
	}

	e.Logger.Info().
		Int64("cycleId", e.CycleId).
		Str("strategy", strg.Name()).
		Object("bestArbirageOpportunity", opportunity).
		Msg("a profitable opportunity found.")

	err = e.Execute(strg.Name(), opportunity.NobitexOrderCandidate)
	if err != nil {
		return fmt.Errorf("failed to execute strategy %s. %w", strg.Name(), err)
	}

	err = e.Execute(strg.Name(), opportunity.UniV3OrderCandidate)
	if err != nil {
		return fmt.Errorf("failed to execute strategy %s. %w", strg.Name(), err)
	}

	_, err = e.Store.CreateNewTrade(context.Background(), e.PairId, e.OrderId, e.TransactionId)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&opportunity)
	if err != nil {
		return err
	}

	err = e.Store.CreateArbitrageOpporchunity(context.Background(), data)
	if err != nil {
		return err
	}

	strg.Teardown()

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
	e.CycleId = cycleId
}

func (e *Executor) ExecuteUniswap(strategyName string, orderCandidate strategy.OrderCandidate) error {
	var (
		txID   int64
		txHash common.Hash
	)
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Object("orderCandidate", orderCandidate).Msg("executing uniswap order candidate")

	tx, err := e.DexTrader.Trade(orderCandidate.Source(), orderCandidate.Destination(), e.UniswapFee, orderCandidate.In, orderCandidate.MinOut)
	if err != nil {
		return err
	}
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("transaction sent to uniswap.")
	txID, err = e.Store.CreateTransaction(context.Background(), tx, e.DexTrader.GetExecutorAddress())
	if err != nil {
		return err
	}
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Int64("transaction_id", txID).Msg("transaction created in database.")

	txHash = tx.Hash()
	e.TransactionId = txID

	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("waiting for transaction receipt.")
	rec, header, err := e.Indxer.WaitForReceipt(context.Background(), txHash)
	if err != nil {
		return err
	}
	txUpdate := transaction.UpdatedFields{
		GasUsed:     rec.GasUsed,
		BlockNumber: rec.BlockNumber.Int64(),
		Timestamp:   time.Unix(int64(header.Time), 0),
		Status:      transaction.CastFromReceiptStatus(rec.Status),
	}
	err = e.Store.UpdateTransaction(context.Background(), txID, txUpdate)
	if err != nil {
		return err
	}
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Str("transaction", tx.Hash().String()).Msg("transaction receipt received.")
	return nil
}

func (e *Executor) ExecuteNobitex(strategyName string, orderCandidate strategy.OrderCandidate) error {
	var (
		nobitexOrderId string
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

	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Object("order", o).Msg("executing nobitex order candidate")

	id, t, err := e.Nobitex.PlaceOrder(o)
	if err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	if id == "" {
		return fmt.Errorf("invalid order id: %d", id)
	}

	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Str("nobitexOrderId", id).Msg("order placed in nobitex.")

	nobitexOrderId = id
	o.CreatedAt = t
	o.OrderId = nobitexOrderId
	o.Status = order.Open
	orderID, err := e.Store.CreateNewOrder(context.Background(), o)
	orderId = orderID
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	e.Logger.Info().Int64("cycleId", e.CycleId).Str("strategy", strategyName).Int64("orderId", orderId).Msg("order created in database.")

	e.OrderId = orderID

	e.Logger.Info().
		Int64("cycleId", e.CycleId).
		Str("strategy", strategyName).
		Str("nobitexOrderId", nobitexOrderId).Msg("waiting for nobitex order status.")
	ctx, cancel := context.WithTimeout(context.Background(), e.NobitexRetryTimeOut)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			e.Logger.Error().Str("nobitexOrderId", nobitexOrderId).Msg("nobitex order status timeout")
			return ctx.Err()
		default:
			orderUpdate, err := e.Nobitex.OrderStatus(context.Background(), nobitexOrderId)
			if err != nil {
				return err
			}
			if orderUpdate.Status != order.Filled {
				continue
			}
			e.Logger.Info().
				Str("nobitexOrderId", nobitexOrderId).
				Str("status", orderUpdate.Status.String()).Msg("nobitex order status received.")
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
			err = e.Store.UpdateOrder(context.Background(), orderId, update)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
