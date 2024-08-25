package uniswapv3

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zarbanio/market-maker-keeper/abis/uniswapv3_quoter"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/math"
)

type Quoter struct {
	quoterAddress common.Address
	quoter        *uniswapv3_quoter.Uniswapv3Quoter
	eth           *ethclient.Client
}

func NewQuoter(eth *ethclient.Client, quoterAddr common.Address) *Quoter {
	quoter, err := uniswapv3_quoter.NewUniswapv3Quoter(quoterAddr, eth)
	if err != nil {
		log.Fatal(err)
	}

	return &Quoter{
		quoterAddress: quoterAddr,
		quoter:        quoter,
		eth:           eth,
	}
}

func (q *Quoter) GetSwapOutputWithExactInput(ctx context.Context, tokenIn, tokenOut domain.Token, fee domain.UniswapFee, amountIn decimal.Decimal) (decimal.Decimal, error) {
	quoterABI, _ := uniswapv3_quoter.Uniswapv3QuoterMetaData.GetAbi()

	amountInWithDecimals, err := math.DecimalToBigInt(amountIn, int32(tokenIn.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}

	data, err := quoterABI.Pack(
		"quoteExactInputSingle",
		tokenIn.Address(),
		tokenOut.Address(),
		fee.BigInt(),
		amountInWithDecimals,
		big.NewInt(0),
	)
	if err != nil {
		return decimal.Zero, err
	}

	output, err := q.eth.CallContract(ctx, ethereum.CallMsg{
		To:   &q.quoterAddress,
		Data: data,
	}, nil)
	if err != nil {
		return decimal.Zero, err
	}

	amountOutWithDecimals := new(big.Int).SetBytes(output)
	amountOut, err := math.BigIntToDecimal(amountOutWithDecimals, int32(tokenOut.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}

	return amountOut, err
}

func (q *Quoter) GetSwapInputWithExactOutput(ctx context.Context, tokenIn, tokenOut domain.Token, fee domain.UniswapFee, amountOut decimal.Decimal) (decimal.Decimal, error) {
	quoterABI, _ := uniswapv3_quoter.Uniswapv3QuoterMetaData.GetAbi()

	amountOutWithDecimals, err := math.DecimalToBigInt(amountOut, int32(tokenOut.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}

	data, err := quoterABI.Pack(
		"quoteExactOutputSingle",
		tokenIn.Address(),
		tokenOut.Address(),
		fee.BigInt(),
		amountOutWithDecimals,
		big.NewInt(0),
	)
	if err != nil {
		return decimal.Zero, err
	}

	output, err := q.eth.CallContract(ctx, ethereum.CallMsg{
		To:   &q.quoterAddress,
		Data: data,
	}, nil)
	if err != nil {
		return decimal.Zero, err
	}

	amountInWithDecimals := new(big.Int).SetBytes(output)
	amountIn, err := math.BigIntToDecimal(amountInWithDecimals, int32(tokenIn.Decimals()))
	if err != nil {
		return decimal.Zero, err
	}

	return amountIn, err
}
