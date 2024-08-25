package nobitex

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
)

type MockExchange struct {
	client            *httpClient
	OrderBooks        map[string]orderbook.OrderBook
	Balance           map[symbol.Symbol]decimal.Decimal
	Order             []order.Order
	Trades            []trade.Trade
	minimumOrderToman decimal.Decimal
	updateTicker      time.Duration
}

func NewMockExchange(baseUrl string, timeout time.Duration, minimumOrderToman decimal.Decimal, orderStatusInterval time.Duration, updateInterval time.Duration) domain.Exchange {
	mockExchange := &MockExchange{
		client:            newClientHttp(baseUrl, "", timeout),
		OrderBooks:        make(map[string]orderbook.OrderBook),
		Order:             make([]order.Order, 100),
		Trades:            make([]trade.Trade, 100),
		Balance:           make(map[symbol.Symbol]decimal.Decimal),
		minimumOrderToman: minimumOrderToman,
		updateTicker:      10,
	}

	go mockExchange.updateOrderBook()

	return mockExchange
}

func (e *MockExchange) updateOrderBook() {
	ticker := time.NewTicker(e.updateTicker * time.Second)
	for {
		res, err := e.client.orderBook(string(symbol.USDT), string(symbol.IRT))
		if err == nil {
			mapKey := fmt.Sprintf("%s:%s", symbol.USDT, symbol.IRT)
			e.OrderBooks[mapKey] = res.toOrderBook()
		}

		<-ticker.C
	}
}

func (e *MockExchange) Fees(tradeType trade.Type) decimal.Decimal {
	if tradeType == trade.Maker {
		return decimal.NewFromFloat(0.002)
	} else if tradeType == trade.Taker {
		return decimal.NewFromFloat(0.0025)
	}
	return decimal.Zero
}

func (e *MockExchange) Decimals(s symbol.Symbol) int32 {
	switch s {
	case symbol.DAI:
		return 2
	case symbol.RLS, symbol.IRT, symbol.PersianRial:
		return 0
	default:
		return 0
	}
}

func (e *MockExchange) OrderBook(src, dst symbol.Symbol) (orderbook.OrderBook, error) {
	// Implement a simple order book with random prices for demonstration purposes.
	srcDes := fmt.Sprintf("%s:%s", src, dst)
	return e.OrderBooks[srcDes], nil
}

func (e *MockExchange) ExchangeRate(src, dst symbol.Symbol) (decimal.Decimal, error) {
	// Simulate exchange rate based on order book prices.
	response, err := e.client.exchangeRate(src.LowerCaseString(), dst.LowerCaseString())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("fetch data failed from api. %w", err)
	}
	pairName := fmt.Sprintf(`%s-%s`, strings.ToLower(src.String()), strings.ToLower(dst.String()))

	// Get the JSON-encoded string from the response for the given pair name
	jsonString, ok := response.Stats[pairName]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("response does not contain expected JSON-encoded string")
	}

	latestFloat, err := decimal.NewFromString(jsonString.Latest)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return latestFloat, nil
}

func (e *MockExchange) CancelOrder(order order.Order) error {
	// Implement order cancellation logic (e.g., remove the order from active orders).
	for i, o := range e.Order {
		if o.Id == order.Id {
			e.Order = append(e.Order[:i], e.Order[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found for cancellation")
}

func (m *MockExchange) SetBalance(symbol symbol.Symbol, balance decimal.Decimal) {
	m.Balance[symbol] = balance
}

func (e *MockExchange) Balances() ([]domain.Balance, error) {
	// Return actual account balances.
	var balances []domain.Balance
	for symbol, balance := range e.Balance {
		balances = append(balances, domain.Balance{Symbol: symbol, Balance: balance})
	}
	return balances, nil
}

func (e *MockExchange) PlaceOrder(o order.Order) (int64, time.Time, error) {
	// Simulate placing an order with possible fill and update balances.
	// In a real exchange, you would send the order to the exchange API.
	o.Id = int64(rand.Intn(1000) + 1) // Generate a random order ID for demonstration.
	o.CreatedAt = time.Now()

	if o.Side == order.Buy {
		// Check if you have enough balance to place a buy order.
		balance, exists := e.Balance[o.SrcCurrency]
		if !exists || balance.LessThan(o.TotalPrice) {
			return 0, time.Time{}, fmt.Errorf("insufficient balance for buy order")
		}

		// Update balances after placing a buy order.
		e.Balance[o.SrcCurrency] = balance.Sub(o.TotalPrice)
		e.Balance[o.DstCurrency] = e.Balance[o.DstCurrency].Add(o.Amount)
	} else if o.Side == order.Sell {
		// Check if you have enough balance to place a sell order.
		balance, exists := e.Balance[o.SrcCurrency]
		if !exists || balance.LessThan(o.Amount) {
			return 0, time.Time{}, fmt.Errorf("insufficient balance for sell order")
		}

		// Update balances after placing a sell order.
		e.Balance[o.SrcCurrency] = balance.Sub(o.Amount)
		e.Balance[o.DstCurrency] = e.Balance[o.DstCurrency].Add(o.TotalPrice)
	}
	o.Status = order.Filled
	e.Order = append(e.Order, o) // Add the order to the active orders.
	return o.Id, o.CreatedAt, nil
}

func (e *MockExchange) OrderStatus(ctx context.Context, id int64) (order.Order, error) {
	// Implement order status retrieval logic (e.g., check the order's status with the exchange API).
	for _, o := range e.Order {
		if o.Id == id {
			fmt.Println("=============", o.Status)
			return o, nil
		}
	}
	return order.Order{}, fmt.Errorf("order not found")
}

func (e *MockExchange) Orders(side order.Side, src, dst symbol.Symbol, state order.State) ([]order.Order, error) {
	// Implement order retrieval based on side, srcCurrency, dstCurrency, and state.
	var filteredOrders []order.Order
	for _, o := range e.Order {
		if o.Side == side && o.SrcCurrency == src && o.DstCurrency == dst && o.Status == state {
			filteredOrders = append(filteredOrders, o)
		}
	}
	return filteredOrders, nil
}

func (e *MockExchange) UpdateOrder(id string, state order.State) (order.State, error) {
	// Implement order state update logic (e.g., send an update request to the exchange API).
	for i, o := range e.Order {
		if fmt.Sprintf("%d", o.Id) == id {
			e.Order[i].Status = state
			return state, nil
		}
	}
	return order.Filled, nil
}
func (e *MockExchange) RecentTrades(src, dst symbol.Symbol) ([]trade.Trade, error) {
	return e.Trades, nil
}

func (e *MockExchange) MinimumOrderToman() decimal.Decimal {
	return e.minimumOrderToman
}
