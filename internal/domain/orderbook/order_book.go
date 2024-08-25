package orderbook

import (
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
)

type OrderBook struct {
	Asks []order.Order
	Bids []order.Order
}

func GetOrderBookPrice(book OrderBook, side order.Side, qty decimal.Decimal) decimal.Decimal {
	switch side {
	case order.Buy:
		return GetOrderBookSidePrice(book.Asks, qty)
	case order.Sell:
		return GetOrderBookSidePrice(book.Bids, qty)
	}
	return decimal.Zero
}

func GetOrderBookSidePrice(orders []order.Order, qty decimal.Decimal) decimal.Decimal {
	sumPrice := decimal.Zero
	totalAmount := decimal.Zero

	for _, o := range orders {
		price := o.Price
		amount := o.Amount

		if o.Amount.Add(totalAmount).LessThan(qty) {
			sumPrice = sumPrice.Add(price.Mul(amount)) // sumPrice += price * amount
			totalAmount = totalAmount.Add(amount)      // totalAmount += amount;
		} else {
			sumPrice = sumPrice.Add(price.Mul(qty.Sub(totalAmount)))
			totalAmount = qty
			break
		}
	}
	if !totalAmount.Equal(qty) {
		return decimal.Zero
	}
	return sumPrice.Div(totalAmount)
}
