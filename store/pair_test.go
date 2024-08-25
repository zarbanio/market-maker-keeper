package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zarbanio/market-maker-keeper/internal/domain/pair"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

func TestPairOperations(t *testing.T) {
	psql := NewPostgres("localhost", 5432, "postgres", "postgres", "market_maker_test")
	// Initialize your Postgres connection here

	// Create a new newPair
	newPair := pair.Pair{
		BaseAsset:  symbol.DAI,
		QuoteAsset: symbol.BTC,
	}

	// Test CreatePair
	pairID, err := psql.CreatePairIfNotExist(context.Background(), &newPair)
	require.NoError(t, err)
	require.NotZero(t, pairID)

	// Test GetPairById
	createdPair, err := psql.GetPairById(context.Background(), pairID)
	require.NoError(t, err)
	require.Equal(t, newPair.BaseAsset, createdPair.BaseAsset)
	require.Equal(t, newPair.QuoteAsset, createdPair.QuoteAsset)
	// Test GetPairList
	pairList, err := psql.GetPairList(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, pairList)

	// Verify that the created newPair exists in the newPair list
	var found bool
	for _, p := range pairList {
		if p.Id == pairID {
			found = true
			break
		}
	}
	require.True(t, found)
}
