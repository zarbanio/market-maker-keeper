package keystore

import (
	"crypto/ecdsa"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyStore struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func New(key string) (*KeyStore, error) {
	ok := strings.HasPrefix(key, "0x")
	if !ok {
		key = "0x" + key
	}
	privateKeyBytes, err := hexutil.Decode(key)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert bytes to private key: %v", err)
	}

	wallet := &KeyStore{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}

	return wallet, nil
}

func (s KeyStore) PrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}

func (s KeyStore) PrivateKeyBytes() []byte {
	return crypto.FromECDSA(s.PrivateKey())
}

func (s KeyStore) PrivateKeyHex() string {
	return hexutil.Encode(s.PrivateKeyBytes())[2:]
}

func (s KeyStore) PublicKey() *ecdsa.PublicKey {
	return s.publicKey
}

func (s KeyStore) PublicKeyBytes() []byte {
	return crypto.FromECDSAPub(s.PublicKey())
}

func (s KeyStore) PublicKeyHex() string {
	return hexutil.Encode(s.PublicKeyBytes())[4:]
}

func (s KeyStore) Address() common.Address {
	return crypto.PubkeyToAddress(*s.publicKey)
}
