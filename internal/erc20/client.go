package erc20

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/abis/IERC20"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/math"
)

type Client struct {
	domain.Token
	ierc20 *IERC20.IERC20
	eth    *ethclient.Client
}

func New(eth *ethclient.Client, token domain.Token) *Client {
	ierc20, err := IERC20.NewIERC20(token.Address(), eth)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Token:  token,
		eth:    eth,
		ierc20: ierc20,
	}
}

func (c *Client) BalanceOf(ctx context.Context, owner common.Address) (decimal.Decimal, error) {
	balance, err := c.ierc20.BalanceOf(&bind.CallOpts{Context: ctx}, owner)
	if err != nil {
		return decimal.Zero, err
	}
	decimals := c.Decimals()

	return math.Normalize(balance, int64(decimals)), nil
}
