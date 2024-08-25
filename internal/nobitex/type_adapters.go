package nobitex

import (
	"strings"

	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

func (r orderBookResponse) toOrderBook() orderbook.OrderBook {
	book := orderbook.OrderBook{
		Asks: make([]order.Order, 0, len(r.Asks)),
		Bids: make([]order.Order, 0, len(r.Bids)),
	}
	for _, ask := range r.Asks {
		book.Asks = append(book.Asks, order.Order{
			Price:  NewDecimalFromString(ask[0]),
			Amount: NewDecimalFromString(ask[1]),
		})
	}
	for _, bid := range r.Bids {
		book.Bids = append(book.Bids, order.Order{
			Price:  NewDecimalFromString(bid[0]),
			Amount: NewDecimalFromString(bid[1]),
		})
	}
	return book
}

func (r walletBalancesResponse) toBalances() []domain.Balance {
	arr := make([]domain.Balance, 0, len(r.Wallets))
	for _, w := range r.Wallets {
		arr = append(arr, domain.Balance{
			Symbol:  symbol.Symbol(strings.ToUpper(w.Currency)),
			Balance: NewDecimalFromString(w.ActiveBalance),
		})
	}
	return arr
}

func (r OrderStatusResponse) toOrder() (order.Order, error) {
	srcCurrency, err := symbol.FromString(r.Order.SrcCurrency)
	if err != nil {
		return order.Order{}, err
	}
	dstCurrency, err := symbol.FromString(r.Order.DstCurrency)
	if err != nil {
		return order.Order{}, err
	}
	side, err := order.SideFromString(r.Order.Type)
	if err != nil {
		return order.Order{}, err
	}

	return order.Order{
		OrderId:         int64(r.Order.Id),
		Side:            side,
		Price:           NewDecimalFromString(r.Order.AveragePrice),
		SrcCurrency:     srcCurrency,
		DstCurrency:     dstCurrency,
		TotalPrice:      NewDecimalFromString(r.Order.TotalPrice),
		TotalOrderPrice: NewDecimalFromString(r.Order.TotalOrderPrice),
		Fee:             NewDecimalFromString(r.Order.Fee),
		Amount:          NewDecimalFromString(r.Order.Amount),
		UnmatchedAmount: NewDecimalFromString(r.Order.UnmatchedAmount),
		MatchedAmount:   NewDecimalFromString(r.Order.MatchedAmount),
		Partial:         r.Order.Partial,
		Status:          toOrderState(r.Order.Status, r.Order.Partial),
		CreatedAt:       r.Order.CreatedAt,
	}, nil
}

func toOrderState(s string, partial bool) order.State {
	switch s {
	case "open", "Open":
		return order.Open
	case "done", "Done":
		if partial {
			return order.PartiallyFilled
		}
		return order.Filled
	case "closed", "Closed":
		return order.Canceled
	}
	return 1000 // TODO
}

func toOrderSide(side string) order.Side {
	switch side {
	case "buy":
		return order.Buy
	case "sell":
		return order.Sell
	}
	return 1000 // TODO
}

func toOrderExecution(t string) order.Execution {
	switch t {
	case "limit":
		return order.LimitExecution
	case "market":
		return order.MarketExecution
	case "stop_market":
		return order.StopMarketExecution
	case "stop_limit":
		return order.StopLimitExecution
	}
	return 1000 // TODO
}
