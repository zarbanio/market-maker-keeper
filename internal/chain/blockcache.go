package chain

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockCache struct {
	client      *ethclient.Client
	blockCache  map[uint64]*types.Block
	headerCache map[uint64]*types.Header
	cacheMux    sync.Mutex
}

func NewBlockCache(eth *ethclient.Client) *BlockCache {
	return &BlockCache{
		client:      eth,
		blockCache:  make(map[uint64]*types.Block),
		headerCache: make(map[uint64]*types.Header),
	}
}

func (bc *BlockCache) GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	bc.cacheMux.Lock()
	defer bc.cacheMux.Unlock()

	block, ok := bc.blockCache[number]
	if ok {
		return block, nil
	}

	block, err := bc.client.BlockByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return nil, err
	}

	bc.blockCache[number] = block
	return block, nil
}

func (bc *BlockCache) GetHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error) {
	bc.cacheMux.Lock()
	defer bc.cacheMux.Unlock()

	header, ok := bc.headerCache[number]
	if ok {
		return header, nil
	}

	header, err := bc.client.HeaderByNumber(ctx, big.NewInt(int64(number)))
	if err != nil {
		return nil, err
	}

	bc.headerCache[number] = header
	return header, nil
}
