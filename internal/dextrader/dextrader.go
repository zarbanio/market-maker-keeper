package dextrader

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/abis/dex_trader"
	"github.com/zarbanio/market-maker-keeper/internal/chain"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
	"github.com/zarbanio/market-maker-keeper/internal/erc20"
	"github.com/zarbanio/market-maker-keeper/internal/keystore"
	"github.com/zarbanio/market-maker-keeper/math"
)

type client struct {
	eth              *ethclient.Client
	dexTraderAddress common.Address
	dexTrader        *dex_trader.DexTrader
	executor         *keystore.KeyStore
	erc20Clients     map[symbol.Symbol]*erc20.Client
}

func newClient(eth *ethclient.Client, dexTraderAddr common.Address, executorWallet *keystore.KeyStore) *client {
	trader, err := dex_trader.NewDexTrader(dexTraderAddr, eth)
	if err != nil {
		log.Fatal(err)
	}
	return &client{
		eth:              eth,
		dexTraderAddress: dexTraderAddr,
		dexTrader:        trader,
		executor:         executorWallet,
		erc20Clients:     make(map[symbol.Symbol]*erc20.Client),
	}
}

func (c *client) getERC20Client(sym symbol.Symbol) (*erc20.Client, error) {
	if value, ok := c.erc20Clients[sym]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("erc20 client for token %s doesn't exist", sym)
}

func (c *client) AddERC20Client(client *erc20.Client) {
	c.erc20Clients[client.Symbol()] = client
}

func (c *client) GetTokenBalances(ctx context.Context) (map[symbol.Symbol]domain.Balance, error) {
	balances := make(map[symbol.Symbol]domain.Balance)

	for _, client := range c.erc20Clients {
		balance, err := client.BalanceOf(ctx, c.dexTraderAddress)
		if err != nil {
			return nil, err
		}

		balances[client.Symbol()] = domain.Balance{
			Symbol:  client.Symbol(),
			Balance: balance,
		}
	}

	return balances, nil
}

func (c *client) GetNativeBalance(ctx context.Context) (decimal.Decimal, error) {
	balance, err := c.eth.BalanceAt(ctx, c.dexTraderAddress, nil)
	if err != nil {
		return decimal.Zero, err
	}

	return math.Normalize(balance, 18), nil
}

func (c *client) GetExecutorNativeBalance(ctx context.Context) (decimal.Decimal, error) {
	balance, err := c.eth.BalanceAt(ctx, c.executor.Address(), nil)
	if err != nil {
		return decimal.Zero, err
	}

	return math.Normalize(balance, 18), nil
}

func (c *client) GetExecutorAddress() common.Address {
	return c.executor.Address()
}

func (c *client) Trade(ctx context.Context, token0, token1 domain.Token, poolFee *big.Int, amountsIn, amountOutMinimum decimal.Decimal) (*types.Transaction, error) {
	gasPrice, err := c.eth.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := chain.GetAccountAuth(ctx, c.eth, c.executor.PrivateKey())
	if err != nil {
		return nil, err
	}

	auth.GasPrice = gasPrice
	tx, err := c.dexTrader.Trade(auth, token0.Address(), token1.Address(), poolFee, math.Denormalize(amountsIn, token0.Decimals()), math.Denormalize(amountOutMinimum, token1.Decimals()))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *client) EstimateDexTradeGasFee(token0, token1 domain.Token, poolFee *big.Int, amountIn, amountOutMinimum decimal.Decimal) (decimal.Decimal, error) {
	amountInWithDecimals, err := math.DecimalToBigInt(amountIn, int32(token0.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}
	amountOutMinimumWithDecimals, err := math.DecimalToBigInt(amountOutMinimum, int32(token1.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}

	// Get gas price
	gasPrice, err := c.eth.SuggestGasPrice(context.Background())
	if err != nil {
		return decimal.Zero, err
	}

	// Estimate the gas for the transaction
	dexTraderABI, _ := dex_trader.DexTraderMetaData.GetAbi()
	data, err := dexTraderABI.Pack("trade", token0.Address(), token1.Address(), poolFee, amountInWithDecimals, amountOutMinimumWithDecimals)
	if err != nil {
		return decimal.Zero, err
	}

	gas, err := c.eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From: c.executor.Address(),
		To:   &c.dexTraderAddress,
		Data: data,
	})
	if err != nil {
		return decimal.Zero, err
	}

	bigIntGas := new(big.Int).SetUint64(gas)
	bigIntGasFee := new(big.Int).Mul(bigIntGas, gasPrice)

	gasFee, err := math.BigIntToDecimal(bigIntGasFee, int32(token0.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}
	return gasFee, nil
}
