package nobitex

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/orderbook"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/domain/trade"
)

type nobitex struct {
	client              *httpClient
	minimumOrderToman   decimal.Decimal
	orderStatusInterval time.Duration
}

func New(baseUrl, apikey string, timeout time.Duration, minimumOrderToman decimal.Decimal, orderStatusInterval time.Duration) domain.Exchange {
	return &nobitex{
		client:              newClientHttp(baseUrl, apikey, timeout),
		minimumOrderToman:   minimumOrderToman,
		orderStatusInterval: orderStatusInterval,
	}
}

func (n nobitex) Fees(tradeType trade.Type) decimal.Decimal {
	if tradeType == trade.Maker {
		return decimal.NewFromFloat(0.002)
	} else if tradeType == trade.Taker {
		return decimal.NewFromFloat(0.0025)
	}
	return decimal.Zero
}

func (n nobitex) Decimals(s symbol.Symbol) int32 {
	switch s {
	case symbol.DAI:
		return 2
	case symbol.RLS, symbol.IRT, symbol.PersianRial:
		return 0
	default:
		return 0
	}
}

func (n nobitex) OrderBook(src, dst symbol.Symbol) (orderbook.OrderBook, error) {
	res, err := n.client.orderBook(src.String(), dst.String())
	if err != nil {
		return orderbook.OrderBook{}, err
	}
	return res.toOrderBook(), nil
}

func (n nobitex) ExchangeRate(src, dst symbol.Symbol) (price decimal.Decimal, err error) {
	response, err := n.client.exchangeRate(src.LowerCaseString(), dst.LowerCaseString())
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
func (n nobitex) SetBalance(symbol symbol.Symbol, balance decimal.Decimal) {
	panic("not implemented")
}

func (n nobitex) Balances() ([]domain.Balance, error) {
	res, err := n.client.walletBalances()
	if err != nil {
		return nil, err
	}
	return res.toBalances(), nil
}

func (n nobitex) PlaceOrder(order order.Order) (string, time.Time, error) {
	if order.DstCurrency == symbol.TMN {
		order.DstCurrency = symbol.RLS
	}
	price, err := n.ExchangeRate(order.SrcCurrency, order.DstCurrency)
	if err != nil {
		return "", time.Time{}, err
	}
	order.Price = price
	res, err := n.client.placeOrder(order)
	if err != nil {
		return "", time.Time{}, err
	}
	return fmt.Sprintf("%d", res.Order.Id), res.Order.CreatedAt, nil
}

func (n nobitex) OrderStatus(ctx context.Context, id string) (order.Order, error) {
	ticker := time.NewTicker(n.orderStatusInterval)
	defer ticker.Stop()

	var _order order.Order
	for {
		select {
		case <-ctx.Done():
			return order.Order{}, ctx.Err()
		case <-ticker.C:
			nobitexOrderId, _ := strconv.ParseInt(id, 10, 64)
			res, err := n.client.orderStatus(nobitexOrderId)
			if err != nil {
				continue
			}
			_order, err = res.toOrder()
			if err == nil {
				return _order, nil
			}
		}
	}
}

func (n nobitex) Orders(side order.Side, src, dst symbol.Symbol, state order.State) ([]order.Order, error) {
	res, err := n.client.orders(toNobitex(side), toNobitex(src), toNobitex(dst), toNobitex(state))
	if err != nil {
		return nil, err
	}
	arr := make([]order.Order, 0, len(res.Orders))
	for _, ord := range res.Orders {
		srcSym, err := symbol.FromString(ord.SrcCurrency)
		if err != nil {
			return nil, err
		}
		dstSym, err := symbol.FromString(ord.DstCurrency)
		if err != nil {
			return nil, err
		}
		exc, err := order.ExecutionFromString(ord.Execution)
		if err != nil {
			return nil, err
		}
		price, err := decimal.NewFromString(ord.Price)
		if err != nil {
			return nil, err
		}
		status, err := order.StateFromString(ord.Status)
		if err != nil {
			return nil, err
		}
		arr = append(arr, order.Order{
			Side:        toOrderSide(ord.Type),
			Execution:   exc,
			SrcCurrency: srcSym,
			DstCurrency: dstSym,
			Price:       NewDecimalFromString(ord.Price),
			Amount:      NewDecimalFromString(ord.Amount),
			StopPrice:   decimal.Zero,
			Fee:         NewDecimalFromString(ord.Fee),
			OrderId:     fmt.Sprintf("%s", ord.Id),
			TotalPrice:  price,
			Status:      status,
			Account:     "", // TODO add account,
			CreatedAt:   time.Now().UTC(),
		})
	}
	return arr, nil
}

func (n nobitex) UpdateOrder(id string, state order.State) (order.State, error) {
	res, err := n.client.updateOrder(id, toNobitex(state))
	if err != nil {
		return 1000, err // TODO
	}
	if res.Status != "ok" {
		return 1000, fmt.Errorf("failed to update order %s state", id) // TODO
	}
	return toOrderState(res.UpdatedStatus, false), nil
}

func (n nobitex) CancelOrder(order order.Order) error {
	res, err := n.client.cancelOrder(order)
	if err != nil {
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("failed to cancel order %s state", strconv.Itoa(int(order.Id))) // TODO
	}
	return nil
}

func (n nobitex) RecentTrades(src, dst symbol.Symbol) ([]trade.Trade, error) {
	res, err := n.client.trades(toNobitex(src), toNobitex(dst), 0)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("failed to get recent trades")
	}
	arr := make([]trade.Trade, 0, len(res.Trades))
	// for _, t := range res.Trades {
	// 	srcSym, err := symbol.FromString(t.SrcCurrency)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	dstSym, err := symbol.FromString(t.DstCurrency)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// arr = append(arr, interface{})
	// }
	return arr, nil
}

func (n nobitex) MinimumOrderToman() decimal.Decimal {
	return n.minimumOrderToman
}
