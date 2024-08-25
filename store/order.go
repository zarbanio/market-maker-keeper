package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type orderModel struct {
	id              int64
	orderId         int64
	execution       string
	side            string
	srcCurrency     string
	dstCurrency     string
	price           string
	amount          string
	totalPrice      string
	totalOrderPrice string
	stopPrice       string
	matchedAmount   string
	unmatchedAmount string
	status          string
	partial         bool
	fee             string
	feeCurrency     string
	account         string
	createdAt       time.Time
}

func (o orderModel) toDomain() (error, *order.Order) {
	srcCurrency, err := symbol.FromString(o.srcCurrency)
	if err != nil {
		return err, nil
	}
	dstCurrency, err := symbol.FromString(o.dstCurrency)
	if err != nil {
		return err, nil
	}
	price, err := decimal.NewFromString(o.price)
	if err != nil {
		return err, nil
	}
	amount, err := decimal.NewFromString(o.amount)
	if err != nil {
		return err, nil
	}
	totalPrice, err := decimal.NewFromString(o.totalPrice)
	if err != nil {
		return err, nil
	}
	totalOrderPrice, err := decimal.NewFromString(o.totalOrderPrice)
	if err != nil {
		return err, nil
	}
	stopPrice, err := decimal.NewFromString(o.stopPrice)
	if err != nil {
		return err, nil
	}
	fee, err := decimal.NewFromString(o.fee)
	if err != nil {
		return err, nil
	}
	feeCurrency, err := symbol.FromString(o.feeCurrency)
	if err != nil {
		return err, nil
	}

	execution, err := order.ExecutionFromString(o.execution)
	if err != nil {
		return err, nil
	}
	side, err := order.SideFromString(o.side)
	if err != nil {
		return err, nil
	}
	um, err := decimal.NewFromString(o.unmatchedAmount)
	if err != nil {
		// handle error
	}
	ma, err := decimal.NewFromString(o.matchedAmount)
	if err != nil {
		// handle error
	}

	state, err := order.StateFromString(o.status)
	if err != nil {
		return err, nil
	}

	return nil, &order.Order{
		Id:              o.id,
		OrderId:         o.orderId,
		Execution:       execution,
		Side:            side,
		SrcCurrency:     srcCurrency,
		DstCurrency:     dstCurrency,
		Price:           price,
		Amount:          amount,
		TotalPrice:      totalPrice,
		TotalOrderPrice: totalOrderPrice,
		StopPrice:       stopPrice,
		Fee:             fee,
		FeeCurrency:     feeCurrency,
		Account:         o.account,
		CreatedAt:       o.createdAt,
		UnmatchedAmount: um,
		MatchedAmount:   ma,
		Partial:         o.partial,
		Status:          state,
	}

}

func (p postgres) CreateNewOrder(ctx context.Context, order order.Order) (int64, error) {

	// prepare insert statement
	stmt := `
        INSERT INTO orders (       
			execution,
			side,
			srcCurrency,
			dstCurrency,
			price,
			amount,
			totalPrice,
			totalOrderPrice,
			stopPrice,
			matchedAmount,
			unmatchedAmount,
			status,
			partial,
			fee,
			feeCurrency,
			account,
			createdAt,
			order_id
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
    `
	var id int64

	err := p.conn.QueryRow(context.Background(), stmt,
		order.Execution,
		order.Side,
		order.SrcCurrency,
		order.DstCurrency,
		order.Price,
		order.Amount,
		order.TotalPrice,
		order.TotalOrderPrice,
		order.StopPrice,
		order.MatchedAmount,
		order.UnmatchedAmount,
		order.Status,
		order.Partial,
		order.Fee,
		order.FeeCurrency,
		order.Account,
		order.CreatedAt.UTC(),
		order.OrderId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}
	return id, nil
}

func (p postgres) UpdateOrder(ctx context.Context, orderID int64, updatedFields order.UpdatedFields) error {
	// Prepare update statement
	stmt := `
	UPDATE orders
	SET 
		matchedAmount = COALESCE($2, matchedAmount),
		unmatchedAmount = COALESCE($3, unmatchedAmount),
		status = COALESCE($4, status),
		fee = COALESCE($5, fee),
		price = COALESCE($6, price),
		totalPrice = COALESCE($7, totalPrice),
		totalOrderPrice = COALESCE($8, totalOrderPrice),
		createdAt = COALESCE($9, createdAt)
	WHERE id = $1
`
	_, err := p.conn.Exec(
		ctx,
		stmt,
		orderID,
		updatedFields.MatchedAmount,
		updatedFields.UnmatchedAmount,
		updatedFields.Status,
		updatedFields.Fee.String(),
		updatedFields.Price.String(),
		updatedFields.TotalPrice.String(),
		updatedFields.TotalOrderPrice.String(),
		updatedFields.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func (p postgres) GetOrderById(ctx context.Context, id int64) (*order.Order, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT *
		FROM orders
		WHERE id = $1
	`, id)

	pm := &orderModel{}

	err := row.Scan(

		&pm.id,
		&pm.orderId,
		&pm.execution,
		&pm.side,
		&pm.srcCurrency,
		&pm.dstCurrency,
		&pm.price,
		&pm.amount,
		&pm.totalPrice,
		&pm.totalOrderPrice,
		&pm.stopPrice,
		&pm.matchedAmount,
		&pm.unmatchedAmount,
		&pm.status,
		&pm.partial,
		&pm.fee,
		&pm.feeCurrency,
		&pm.account,
		&pm.createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	err, order := pm.toDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (p postgres) GetOrderByOrderId(ctx context.Context, orderID int64) (*order.Order, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT *
		FROM orders
		WHERE order_id = $1
	`, orderID)

	pm := &orderModel{}

	err := row.Scan(

		&pm.id,
		&pm.orderId,
		&pm.execution,
		&pm.side,
		&pm.srcCurrency,
		&pm.dstCurrency,
		&pm.price,
		&pm.amount,
		&pm.totalPrice,
		&pm.totalOrderPrice,
		&pm.stopPrice,
		&pm.matchedAmount,
		&pm.unmatchedAmount,
		&pm.status,
		&pm.partial,
		&pm.fee,
		&pm.feeCurrency,
		&pm.account,
		&pm.createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	err, order := pm.toDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}
func (p postgres) GetOrders(ctx context.Context, page int, pageSize int) ([]order.Order, error) {
	// Calculate the offset based on the page and page size
	offset := (page - 1) * pageSize

	rows, err := p.conn.Query(ctx, `
			SELECT *
			FROM orders
			ORDER BY created_at DESC
			OFFSET $1
			LIMIT $2
		`, offset, pageSize)

	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		pm := &orderModel{}

		err = rows.Scan(
			&pm.id,
			&pm.orderId,
			&pm.execution,
			&pm.side,
			&pm.srcCurrency,
			&pm.dstCurrency,
			&pm.price,
			&pm.amount,
			&pm.totalPrice,
			&pm.totalOrderPrice,
			&pm.stopPrice,
			&pm.matchedAmount,
			&pm.unmatchedAmount,
			&pm.status,
			&pm.partial,
			&pm.fee,
			&pm.feeCurrency,
			&pm.account,
			&pm.createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		err, normalizedPm := pm.toDomain()
		if err != nil {
			return nil, err
		}

		orders = append(orders, *normalizedPm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over orders: %w", err)
	}

	return orders, nil
}
