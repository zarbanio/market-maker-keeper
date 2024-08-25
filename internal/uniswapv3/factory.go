package uniswapv3

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zarbanio/market-maker-keeper/abis/uniswapv3_factory"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
)

type Factory struct {
	factory *uniswapv3_factory.Uniswapv3Factory
	eth     *ethclient.Client
}

func NewFactory(eth *ethclient.Client, factoryAddr common.Address) Factory {
	factory, err := uniswapv3_factory.NewUniswapv3Factory(factoryAddr, eth)
	if err != nil {
		log.Fatal(err)
	}
	return Factory{
		factory: factory,
		eth:     eth,
	}
}

func (f *Factory) GetPool(ctx context.Context, token0, token1 common.Address, fee domain.UniswapFee) (common.Address, error) {
	return f.factory.GetPool(&bind.CallOpts{Context: ctx}, token0, token1, fee.BigInt())
}
