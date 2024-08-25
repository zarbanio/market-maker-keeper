package uniswapv3

import (
	"context"
	"log"
	"math/big"

	"github.com/zarbanio/market-maker-keeper/math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zarbanio/market-maker-keeper/abis/uniswapv3_pool"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/mstore"
)

type Pool struct {
	pool          *uniswapv3_pool.Uniswapv3Pool
	eth           *ethclient.Client
	memTokenStore mstore.IToken

	token0 domain.Token
	token1 domain.Token
}

func NewPool(eth *ethclient.Client, poolAddr common.Address, memTokenStore mstore.IToken) Pool {
	pool, err := uniswapv3_pool.NewUniswapv3Pool(poolAddr, eth)
	if err != nil {
		log.Fatal(err)
	}

	// get token0 from pool
	token0Address, err := pool.Token0(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}

	token0, err := memTokenStore.GetTokenByAddress(token0Address)
	if err != nil {
		log.Fatal(err)
	}

	// get token1 from pool
	token1Address, err := pool.Token1(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}

	token1, err := memTokenStore.GetTokenByAddress(token1Address)
	if err != nil {
		log.Fatal(err)
	}

	return Pool{
		pool:          pool,
		eth:           eth,
		memTokenStore: memTokenStore,

		token0: token0,
		token1: token1,
	}
}

func (p *Pool) GetBalances(ctx context.Context) (domain.Balance, domain.Balance, error) {
	liquidity, err := p.pool.Liquidity(&bind.CallOpts{Context: ctx})
	if err != nil {
		return domain.Balance{}, domain.Balance{}, err
	}

	slot0, err := p.pool.Slot0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return domain.Balance{}, domain.Balance{}, err
	}

	amount0 := new(big.Int).Div(new(big.Int).Mul(liquidity, math.FixedPoint96), slot0.SqrtPriceX96) // DAI
	amount1 := new(big.Int).Div(new(big.Int).Mul(liquidity, slot0.SqrtPriceX96), math.FixedPoint96) // ZAR

	token0Balance := domain.Balance{
		Symbol:  p.token0.Symbol(),
		Balance: math.Normalize(amount0, p.token0.Decimals()),
	}

	token1Balance := domain.Balance{
		Symbol:  p.token1.Symbol(),
		Balance: math.Normalize(amount1, p.token1.Decimals()),
	}

	return token0Balance, token1Balance, nil
}
