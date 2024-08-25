package nobitex

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
	"github.com/zarbanio/market-maker-keeper/internal/domain/symbol"
)

func TestOrderBook(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a mock orderbook
		// encode the orderbook as JSON
		mockResponse := orderBookResponse{
			Status:     "success",
			LastUpdate: 1647564000,
			Bids: [][]string{
				{"100", "10"},
				{"100", "10"},
				{"100", "10"},
				{"100", "10"},
			},
			Asks: [][]string{
				{"101", "8"},
				{"101", "8"},
			},
		}

		response, err := json.Marshal(mockResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header
		w.Header().Set("Content-Type", "application/json")

		// write the response
		w.Write(response)
	}))
	defer svr.Close()
	// start the server

	nobitexInstance := New(svr.URL, "Key", 0, decimal.NewFromInt(11), 60)

	orders, err := nobitexInstance.OrderBook("DAI", "USDT")

	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	assert.Equal(t, len(orders.Bids), 4, "expected 4 bids")

	// assert that orders.Asks has a length of 2
	assert.Equal(t, len(orders.Asks), 2, "expected 2 asks")
}
func TestOrderBookServerError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer svr.Close()

	nobitexInstance := New(svr.URL, "w", 0, decimal.NewFromInt(11), 60)

	_, err := nobitexInstance.OrderBook("DAI", "USDT")
	t.Log(err)
	assert.Error(t, err, "Expected an error to be returned")
}

func TestPlaceOrder(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mockOrder := placeOrderResponse{
			Status: "success",
			Order: struct {
				Type            string    `json:"type"`
				SrcCurrency     string    `json:"srcCurrency"`
				DstCurrency     string    `json:"dstCurrency"`
				Price           string    `json:"price"`
				Amount          string    `json:"amount"`
				TotalPrice      string    `json:"totalPrice"`
				MatchedAmount   int       `json:"matchedAmount"`
				UnmatchedAmount string    `json:"unmatchedAmount"`
				Id              int       `json:"id"`
				Status          string    `json:"status"`
				Partial         bool      `json:"partial"`
				Fee             int       `json:"fee"`
				User            string    `json:"user"`
				CreatedAt       time.Time `json:"created_at"`
			}{
				Type:            "buy",
				SrcCurrency:     "USD",
				DstCurrency:     "BTC",
				Price:           "50000",
				Amount:          "1",
				TotalPrice:      "50000",
				MatchedAmount:   1,
				UnmatchedAmount: "0",
				Id:              123,
				Status:          "open",
				Partial:         false,
				Fee:             100,
				User:            "johndoe",
				CreatedAt:       time.Now(),
			},
		}

		response, err := json.Marshal(mockOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header
		w.Header().Set("Content-Type", "application/json")

		// write the response
		w.Write(response)
	}))
	defer svr.Close()
	nobitexInstance := New(svr.URL, "w", time.Duration(10), decimal.NewFromInt(11), 60)
	_, err := nobitexInstance.OrderBook("DAI", "USDT")
	assert.Error(t, err, "Expected an error, but got nil")

}

func TestExchangeRate(t *testing.T) {
	mockData := exchangeRateResponse{
		Status: "success",
		Stats: map[string]struct {
			BestSell  string `json:"bestSell"`
			DayOpen   string `json:"dayOpen"`
			DayHigh   string `json:"dayHigh"`
			BestBuy   string `json:"bestBuy"`
			VolumeSrc string `json:"volumeSrc"`
			DayLow    string `json:"dayLow"`
			Latest    string `json:"latest"`
			VolumeDst string `json:"volumeDst"`
			DayChange string `json:"dayChange"`
			DayClose  string `json:"dayClose"`
			IsClosed  bool   `json:"IsClosed"`
		}{
			"btc-usdt": {
				BestSell:  "50000",
				DayOpen:   "49000",
				DayHigh:   "52000",
				BestBuy:   "49000",
				VolumeSrc: "1000",
				DayLow:    "48000",
				Latest:    "51000",
				VolumeDst: "20000000",
				DayChange: "3%",
				DayClose:  "50500",
				IsClosed:  false,
			},
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response, err := json.Marshal(mockData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header
		w.Header().Set("Content-Type", "application/json")

		// write the response
		w.Write(response)
	}))
	defer svr.Close()
	nobitexInstance := New(svr.URL, "w", 0, decimal.NewFromInt(11), 60)
	lastPrice, err := nobitexInstance.ExchangeRate(symbol.BTC, symbol.USDT)

	if err != nil {
		assert.Error(t, err, "not expecte an error, but got err")
	}

	// Assert the last price
	expectedLastPrice, _ := decimal.NewFromString(mockData.Stats["btc-usdt"].Latest)

	assert.Equal(t, expectedLastPrice.String(), lastPrice.String(), "last price mismatch")

}

func TestPlaceOrderServerError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer svr.Close()

	nobitextInstance := New(svr.URL, "w", 0, decimal.NewFromInt(11), 60)
	newOrder := order.Order{
		Id:        1,
		Side:      order.Buy,
		Price:     NewDecimalFromString("40000"),
		StopPrice: NewDecimalFromString("4000"),
		Fee:       NewDecimalFromString("2.3"),
	}
	_, _, err := nobitextInstance.PlaceOrder(newOrder)
	assert.Error(t, err, "Expected an error to be returned")

}

func TestExchangeRateServerError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer svr.Close()

	nobitexInstance := New(svr.URL, "w", 0, decimal.NewFromInt(11), 60)

	_, err := nobitexInstance.ExchangeRate(symbol.BTC, symbol.USDT)
	if err != nil {
		assert.Error(t, err, "not expect an error, but got err")
	}

	assert.Error(t, err, "Expected an error to be returned")
}

func TestCancelOrder(t *testing.T) {
	cancelOrderRes := cancelOrderResponse{
		Status: "success",
		Order: struct {
			Id              int64  `json:"id"`
			Type            string `json:"type"`
			Execution       string `json:"execution"`
			TradeType       string `json:"tradeType"`
			SrcCurrency     string `json:"srcCurrency"`
			DstCurrency     string `json:"dstCurrency"`
			Price           string `json:"price"`
			Amount          string `json:"amount"`
			Status          string `json:"status"`
			TotalPrice      string `json:"totalPrice"`
			TotalOrderPrice string `json:"totalOrderPrice"`
			MatchedAmount   string `json:"matchedAmount"`
			UnmatchedAmount string `json:"unmatchedAmount"`
			Partial         bool   `json:"partial"`
			Fee             int    `json:"fee"`
			CreatedAt       string `json:"createdAt"`
			AveragePrice    string `json:"averagePrice"`
		}{
			Id:              12345,
			Type:            "limit",
			Execution:       "post-only",
			TradeType:       "sell",
			SrcCurrency:     "BTC",
			DstCurrency:     "USD",
			Price:           "50000",
			Amount:          "0.5",
			Status:          "cancelled",
			TotalPrice:      "25000",
			TotalOrderPrice: "25000",
			MatchedAmount:   "0",
			UnmatchedAmount: "0.5",
			Partial:         false,
			Fee:             10,
			CreatedAt:       "2022-04-28T15:45:00Z",
			AveragePrice:    "",
		},
	}

	orderMock := order.Order{
		Id:        12345,
		Execution: order.LimitExecution,

		Price: NewDecimalFromString("50000"),

		Fee:         NewDecimalFromString("10"),
		Side:        order.Buy,
		SrcCurrency: symbol.DAI,
		DstCurrency: symbol.DAI,
		Amount:      NewDecimalFromString("0.5"),
		StopPrice:   NewDecimalFromString("50000"),
		CreatedAt:   time.Now(),
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response, err := json.Marshal(cancelOrderRes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header
		w.Header().Set("Content-Type", "application/json")

		// write the response
		w.Write(response)
	}))
	defer svr.Close()

	nobitexInstance := New(svr.URL, "w", 0, decimal.NewFromInt(11), 60)

	err := nobitexInstance.CancelOrder(orderMock)

	if err != nil {
		assert.Error(t, err, "not expect an error, but got err")
	}
}
