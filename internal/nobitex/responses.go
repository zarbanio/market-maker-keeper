package nobitex

import (
	"time"
)

type (
	cancelOrderResponse struct {
		Status string `json:"status"`
		Order  struct {
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
		} `json:"order"`
	}
	orderBookResponse struct {
		Status     string     `json:"status"`
		LastUpdate int64      `json:"lastUpdate"`
		Bids       [][]string `json:"bids"`
		Asks       [][]string `json:"asks"`
	}

	walletBalancesResponse struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Wallets []struct {
			Currency       string `json:"currency"`
			ActiveBalance  string `json:"activeBalance"`
			BlockedBalance string `json:"blockedBalance"`
			Balance        string `json:"balance"`
		} `json:"wallets"`
	}

	placeOrderResponse struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Order   struct {
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
		} `json:"order"`
	}

	OrderStatusResponse struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Order   struct {
			UnmatchedAmount string    `json:"unmatchedAmount"`
			Fee             string    `json:"fee"`
			MatchedAmount   string    `json:"matchedAmount"`
			Partial         bool      `json:"partial"`
			Price           string    `json:"price"`
			CreatedAt       time.Time `json:"created_at"`
			User            string    `json:"user"`
			Id              int       `json:"id"`
			SrcCurrency     string    `json:"srcCurrency"`
			TotalPrice      string    `json:"totalPrice"`
			TotalOrderPrice string    `json:"totalOrderPrice"`
			Type            string    `json:"type"`
			DstCurrency     string    `json:"dstCurrency"`
			IsMyOrder       bool      `json:"isMyOrder"`
			Status          string    `json:"status"`
			Amount          string    `json:"amount"`
			AveragePrice    string    `json:"averagePrice"`
		} `json:"order"`
	}

	exchangeRateResponse struct {
		Status string `json:"status"`
		Stats  map[string]struct {
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
		} `json:"Stats"`
	}

	ordersResponse struct {
		Status string `json:"status"`
		Orders []struct {
			Id            int    `json:"id"`
			Type          string `json:"type"`
			Execution     string `json:"execution"`
			Status        string `json:"status"`
			SrcCurrency   string `json:"srcCurrency"`
			DstCurrency   string `json:"dstCurrency"`
			Price         string `json:"price"`
			Amount        string `json:"amount"`
			MatchedAmount string `json:"matchedAmount"`
			AveragePrice  string `json:"averagePrice"`
			Fee           string `json:"fee"`
		} `json:"orders"`
	}

	updateOrderResponse struct {
		Status        string `json:"status"`
		UpdatedStatus string `json:"updatedStatus"`
	}

	tradesResponse struct {
		Status string `json:"status"`
		Trades []struct {
			Id          int       `json:"id"`
			OrderId     int       `json:"orderId"`
			SrcCurrency string    `json:"srcCurrency"`
			DstCurrency string    `json:"dstCurrency"`
			Market      string    `json:"market"`
			Timestamp   time.Time `json:"timestamp"`
			Type        string    `json:"type"`
			Price       string    `json:"price"`
			Amount      string    `json:"amount"`
			Total       string    `json:"total"`
			Fee         string    `json:"fee"`
		} `json:"trades"`
		HasNext bool `json:"hasNext"`
	}
)
