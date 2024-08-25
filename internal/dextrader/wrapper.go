package dextrader

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/erc20"
	"github.com/zarbanio/market-maker-keeper/internal/keystore"
	"github.com/zarbanio/market-maker-keeper/internal/mstore"
)

type Wrapper struct {
	client     *client
	tokenStore mstore.IToken
}

func New(eth *ethclient.Client, dexTraderAddr common.Address, executorWallet *keystore.KeyStore, tokenStore mstore.IToken) *Wrapper {
	return &Wrapper{
		client:     newClient(eth, dexTraderAddr, executorWallet),
		tokenStore: tokenStore,
	}
}

func (w *Wrapper) AddERC20Client(client *erc20.Client) {
	w.client.AddERC20Client(client)
}

func (w *Wrapper) GetTokenBalances(ctx context.Context) (map[symbol.Symbol]domain.Balance, error) {
	return w.client.GetTokenBalances(ctx)
}

func (w *Wrapper) GetNativeBalance(ctx context.Context) (decimal.Decimal, error) {
	return w.client.GetNativeBalance(ctx)
}

func (w *Wrapper) GetExecutorNativeBalance(ctx context.Context) (decimal.Decimal, error) {
	return w.client.GetExecutorNativeBalance(ctx)
}

func (w *Wrapper) GetExecutorAddress() common.Address {
	return w.client.GetExecutorAddress()
}

func (w *Wrapper) Trade(src, dst symbol.Symbol, poolFee domain.UniswapFee, quantity, amountOutMinimum decimal.Decimal) (*types.Transaction, error) {
	token0, err := w.tokenStore.GetTokenBySymbol(src)
	if err != nil {
		return &types.Transaction{}, err
	}
	token1, err := w.tokenStore.GetTokenBySymbol(dst)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := w.client.Trade(context.Background(), token0, token1, poolFee.BigInt(), quantity, amountOutMinimum)
	if err != nil {
		return &types.Transaction{}, err
	}
	return tx, nil
}

func (w *Wrapper) EstimateDexTradeGasFee(token0, token1 domain.Token, poolFee *big.Int, amountIn, amountOutMinimum decimal.Decimal) (decimal.Decimal, error) {
	return w.client.EstimateDexTradeGasFee(token0, token1, poolFee, amountIn, amountOutMinimum)
}
