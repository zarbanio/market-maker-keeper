package store

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

func TestGetOrderById(t *testing.T) {
	// Set up a new Postgres instance for testing
	psql := NewPostgres("localhost", 5432, "postgres", "postgres", "market_maker_test")
	err := psql.Migrate("/migrations")
	require.NoError(t, err)

	// Create a new order and insert it into the database
	order := &order.Order{
		Id:              1,
		OrderId:         123,
		Side:            order.Buy,
		SrcCurrency:     symbol.BTC,
		DstCurrency:     symbol.ZAR,
		Price:           decimal.NewFromInt(35000),
		Amount:          decimal.NewFromInt(1),
		TotalPrice:      decimal.NewFromInt(35000),
		TotalOrderPrice: decimal.NewFromInt(35000),
		Status:          order.Open,
		Account:         "test",
		CreatedAt:       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		Execution:       order.LimitExecution,
		StopPrice:       decimal.NewFromInt(0),
		Fee:             decimal.NewFromInt(0),
		FeeCurrency:     symbol.BTC,
		UnmatchedAmount: decimal.NewFromInt(1),
		MatchedAmount:   decimal.NewFromInt(0),
		Partial:         false,
	}

	newId, err := psql.CreateNewOrder(context.Background(), *order)
	require.NoError(t, err)

	// Call the GetOrderById method to retrieve the order we just inserted
	result, err := psql.GetOrderById(context.Background(), newId)
	require.NoError(t, err)

	// Compare the retrieved order with the original order to ensure they match
	require.Equal(t, result.Id, newId)
	require.Equal(t, order.OrderId, result.OrderId)
	require.Equal(t, order.Side, result.Side)
	require.Equal(t, order.SrcCurrency, result.SrcCurrency)
	require.Equal(t, order.DstCurrency, result.DstCurrency)
	require.Equal(t, order.Price, result.Price)
	require.Equal(t, order.Amount, result.Amount)
	require.Equal(t, order.TotalPrice, result.TotalPrice)
	require.Equal(t, order.TotalOrderPrice, result.TotalOrderPrice)
	require.Equal(t, order.Status, result.Status)
	require.Equal(t, order.Account, result.Account)
	require.Equal(t, order.Fee, result.Fee)
	require.Equal(t, order.FeeCurrency, result.FeeCurrency)
	require.True(t, result.CreatedAt.UTC().Equal(order.CreatedAt.UTC()))
}

func TestUpdateOrder(t *testing.T) {
	// Set up a new Postgres instance for testing
	psql := NewPostgres("localhost", 5432, "postgres", "postgres", "market_maker_test")
	err := psql.Migrate("/migrations")
	require.NoError(t, err)

	// Create a new order and insert it into the database
	newOrder := &order.Order{
		OrderId:         123,
		Side:            order.Buy,
		SrcCurrency:     symbol.BTC,
		DstCurrency:     symbol.ZAR,
		Price:           decimal.NewFromInt(35000),
		Amount:          decimal.NewFromInt(1),
		TotalPrice:      decimal.NewFromInt(35000),
		TotalOrderPrice: decimal.NewFromInt(35000),
		Status:          order.Open,
		Account:         "test",
		CreatedAt:       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		Execution:       order.LimitExecution,
		StopPrice:       decimal.NewFromInt(0),
		Fee:             decimal.NewFromInt(0),
		FeeCurrency:     symbol.BTC,
		UnmatchedAmount: decimal.NewFromInt(1),
		MatchedAmount:   decimal.NewFromInt(0),
		Partial:         false,
	}

	newID, err := psql.CreateNewOrder(context.Background(), *newOrder)
	require.NoError(t, err)
	matchedAmount := decimal.NewFromInt(0)
	unmatchedAmount := decimal.NewFromInt(0)
	status := order.Canceled
	fee := decimal.NewFromInt(10)
	feeCurrency := symbol.BTC
	price := decimal.NewFromInt(10)
	totalPrice := decimal.NewFromInt(100)
	totalOrderPrice := decimal.NewFromInt(100)
	createdAt := time.Now()

	// Call the UpdateOrder method to update the order in the database
	err = psql.UpdateOrder(context.Background(), newID, order.UpdatedFields{
		MatchedAmount:   &matchedAmount,
		UnmatchedAmount: &unmatchedAmount,
		Status:          &status,
		Fee:             &fee,
		Price:           &price,
		TotalPrice:      &totalPrice,
		TotalOrderPrice: &totalOrderPrice,
		CreatedAt:       &createdAt,
	})

	require.NoError(t, err)

	// Retrieve the updated order from the database
	updatedOrder, err := psql.GetOrderById(context.Background(), newID)
	require.NoError(t, err)

	// Assert that the order fields have been updated correctly
	require.Equal(t, matchedAmount, updatedOrder.MatchedAmount)
	require.Equal(t, unmatchedAmount, updatedOrder.UnmatchedAmount)
	require.Equal(t, status, updatedOrder.Status)
	require.Equal(t, fee, updatedOrder.Fee)
	require.Equal(t, feeCurrency, updatedOrder.FeeCurrency.String())
	require.Equal(t, price, updatedOrder.Price)
	require.Equal(t, totalPrice, updatedOrder.TotalPrice)
	require.Equal(t, totalOrderPrice, updatedOrder.TotalOrderPrice)

	//In PostgreSQL, the TIMESTAMP data type has a maximum precision of microseconds (6 decimal places).
	require.Equal(t, createdAt.UTC().Truncate(time.Microsecond), updatedOrder.CreatedAt.Truncate(time.Microsecond))
}
