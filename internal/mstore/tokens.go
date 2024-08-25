package mstore

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zarbanio/market-maker-keeper/internal/domain"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

type IToken interface {
	AddToken(token domain.Token) error
	GetTokenBySymbol(symbol symbol.Symbol) (domain.Token, error)
	GetTokenByAddress(address common.Address) (domain.Token, error)
	GetAllTokens() ([]domain.Token, error)
	RemoveToken(symbol symbol.Symbol) error
}

type memTokenStore struct {
	mu  *sync.Mutex
	mem map[symbol.Symbol]domain.Token
}

func NewMemoryTokenStore() IToken {
	return &memTokenStore{
		mu:  &sync.Mutex{},
		mem: make(map[symbol.Symbol]domain.Token),
	}
}

func (m *memTokenStore) AddToken(token domain.Token) error {
	m.mu.Lock()
	m.mem[token.Symbol()] = token
	m.mu.Unlock()
	return nil
}

func (m *memTokenStore) GetTokenBySymbol(symbol symbol.Symbol) (domain.Token, error) {
	m.mu.Lock()
	v, ok := m.mem[symbol]
	m.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("token %s not found", symbol)
	}
	return v, nil
}

func (m *memTokenStore) GetTokenByAddress(address common.Address) (domain.Token, error) {
	m.mu.Lock()
	for _, token := range m.mem {
		if token.Address() == address {
			m.mu.Unlock()
			return token, nil
		}
	}
	m.mu.Unlock()
	return nil, fmt.Errorf("token with address %s not found", address)
}

func (m *memTokenStore) GetAllTokens() ([]domain.Token, error) {
	m.mu.Lock()
	var tokens []domain.Token
	for _, token := range m.mem {
		tokens = append(tokens, token)
	}
	m.mu.Unlock()
	return tokens, nil
}

func (m *memTokenStore) RemoveToken(symbol symbol.Symbol) error {
	m.mu.Lock()
	delete(m.mem, symbol)
	m.mu.Unlock()
	return nil
}
