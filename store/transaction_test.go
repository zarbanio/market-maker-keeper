package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	randm "math/rand"
	"testing"
	"time"

	"github.com/zarbanio/market-maker-keeper/internal/domain/transaction"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zarbanio/market-maker-keeper/internal/domain/blockchain"
)

func TestCreateTransaction(t *testing.T) {
	psql := NewPostgres("localhost", 5432, "postgres", "postgres", "market_maker_test")

	randm.Seed(time.Now().UnixNano()) // Initialize the random number generator
	nonce := uint64(randm.Int63())
	to := common.HexToAddress("0x...")
	value := big.NewInt(100)
	gasLimit := uint64(200000)
	gasPrice := big.NewInt(1000000000)
	data := []byte("...")

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	from := common.HexToAddress("0x...")

	id, err := psql.CreateTransaction(context.Background(), tx, from)
	require.NoError(t, err)

	_, err = psql.GetTransactionById(context.Background(), id)
	require.NoError(t, err)
}
func createTransaction(nonce uint64, toAddress common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	transaction := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	return transaction
}
func calculateTransactionFee(gasLimit uint64, gasPrice *big.Int) string {
	fee := new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPrice)
	return fee.String()
}
func createDummyTransaction() *types.Transaction {
	// Set the transaction fields
	randomUUID := uuid.New()
	randomAddress := fmt.Sprintf("0x%s", randomUUID)
	toAddress := common.HexToAddress(randomAddress)
	value := big.NewInt(1000000000000000000) // 1 ETH
	gasPrice := big.NewInt(1000000000)       // 1 Gwei
	nonce := uint64(0)
	data := []byte("dummy transaction data")
	gasLimit := uint64(21000) // Set a default gas limit value

	// Create the transaction
	transaction := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	return transaction
}
func NewBlockchain(name string, nativeCurrency string, chainID int64) *blockchain.Blockchain {
	return &blockchain.Blockchain{
		Name:           name,
		NativeCurrency: nativeCurrency,
		ChainId:        chainID,
	}
}
func GenerateRandomAddress() (common.Address, error) {
	// Generate a random private key
	privateKey, err := rand.Int(rand.Reader, new(big.Int).SetBytes([]byte("0123456789abcdef"))) // 16-byte private key
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Pad the private key with zero bytes to make it 20 bytes long
	paddedPrivateKey := make([]byte, 20)
	copy(paddedPrivateKey[20-len(privateKey.Bytes()):], privateKey.Bytes())

	// Get the corresponding public key
	publicKey := paddedPrivateKey

	// Generate the Ethereum address from the last 20 bytes of the public key
	address := common.BytesToAddress(publicKey[len(publicKey)-20:])

	return address, nil
}

func TestUpdateTransaction(t *testing.T) {
	// Set up a new Postgres instance for testing
	p := NewPostgres("localhost", 5432, "postgres", "postgres", "market_maker_test")
	err := p.Migrate("/migrations")
	require.NoError(t, err)
	nonce := uint64(2)

	toAddress, err := GenerateRandomAddress()
	require.NoError(t, err)
	value := big.NewInt(1000000000000000000) // 1 ETH
	gasLimit := uint64(21000)                // Set a default gas limit value
	gasPrice := big.NewInt(1000000000)       // 1 Gwei
	data := []byte("dummy transaction data")

	// Create the transaction
	dummyTransaction := createTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	newID, err := p.CreateTransaction(context.Background(), dummyTransaction, *dummyTransaction.To())
	require.NoError(t, err)

	// Create the updated fields
	updatedFields := transaction.UpdatedFields{
		BlockNumber: 348767384,
		GasUsed:     234000,
		Timestamp:   time.Now(),
		Status:      transaction.Success,
	}

	// Call the UpdateTransaction method to update the transaction in the database
	err = p.UpdateTransaction(context.Background(), newID, updatedFields)
	require.NoError(t, err)

	// Retrieve the updated transaction from the database
	updatedTransaction, err := p.GetTransactionById(context.Background(), newID)
	require.NoError(t, err)

	// Assert that the transaction fields have been updated correctly
	require.Equal(t, updatedFields.Status, updatedTransaction.TransactionStatus)
}
