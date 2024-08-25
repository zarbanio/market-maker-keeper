package store

import "errors"

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrPairNotFound        = errors.New("pair not found")
	ErrBlockchainNotFound  = errors.New("blockchain not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrTradeNotFound       = errors.New("trade not found")
)
