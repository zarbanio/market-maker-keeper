package nobitex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/zarbanio/market-maker-keeper/internal/domain/order"
)

// {
//     "type": "sell",
//     "execution": "market",
//     "srcCurrency": "usdt",
//     "dstCurrency": "rls",
//     "amount": 2.763075,
//     "price": 548440
// }

type httpClient struct {
	baseUrl string
	apiKey  string
	cli     *http.Client
}

func newClientHttp(url, apiKey string, timeout time.Duration) *httpClient {
	return &httpClient{
		baseUrl: url,
		cli:     &http.Client{Timeout: timeout},
		apiKey:  apiKey,
	}
}

func (c *httpClient) orderBook(src, dst string) (orderBookResponse, error) {
	u := fmt.Sprintf("%s/v2/orderbook/%s%s", c.baseUrl, src, dst)
	var v orderBookResponse
	err := c.do(http.MethodGet, u, c.authHeader(), nil, &v)
	if err != nil {
		return orderBookResponse{}, err
	}
	return v, nil
}

func (c *httpClient) walletBalances() (walletBalancesResponse, error) {
	u := c.baseUrl + "/users/wallets/list"
	var v walletBalancesResponse
	err := c.do(http.MethodGet, u, c.authHeader(), nil, &v)
	if err != nil {
		return walletBalancesResponse{}, err
	}
	if v.Status == "failed" {
		return walletBalancesResponse{}, fmt.Errorf(`error code: %s, error message: %s`, v.Code, v.Message)
	}
	return v, nil
}

func (c *httpClient) placeOrder(order order.Order) (placeOrderResponse, error) {
	u := c.baseUrl + "/market/orders/add"
	price, _ := order.Price.Float64()
	amount, _ := order.Amount.Float64()
	body, err := json.Marshal(map[string]interface{}{
		"type":        order.Side.String(),
		"execution":   order.Execution.String(),
		"srcCurrency": order.SrcCurrency.LowerCaseString(),
		"dstCurrency": order.DstCurrency.LowerCaseString(),
		"amount":      amount,
		"price":       price,
	})
	if err != nil {
		return placeOrderResponse{}, err
	}
	var v placeOrderResponse
	err = c.do(http.MethodPost, u, c.authHeader(), bytes.NewReader(body), &v)
	if err != nil {
		return placeOrderResponse{}, err
	}
	if v.Status == "failed" {
		return placeOrderResponse{}, fmt.Errorf(`error code: %s, error message: %s`, v.Code, v.Message)
	}
	return v, nil
}

func (c *httpClient) orderStatus(id int64) (OrderStatusResponse, error) {
	u := c.baseUrl + "/market/orders/status"
	body, _ := json.Marshal(map[string]interface{}{
		"id": id,
	})
	var v OrderStatusResponse
	err := c.do(http.MethodPost, u, c.authHeader(), bytes.NewReader(body), &v)
	if err != nil {
		return OrderStatusResponse{}, err
	}
	if v.Status == "failed" {
		return OrderStatusResponse{}, fmt.Errorf(`error code: %s, error message: %s`, v.Code, v.Message)
	}
	return v, nil
}

func (c *httpClient) exchangeRate(srcCurrency string, dstCurrency string) (exchangeRateResponse, error) {
	u := c.baseUrl + "/market/stats" + "?srcCurrency=" + srcCurrency + "&dstCurrency=" + dstCurrency
	var v exchangeRateResponse
	err := c.do(http.MethodGet, u, c.authHeader(), nil, &v)
	if err != nil {
		return v, err
	}
	return v, nil
}

func (c *httpClient) orders(side, src, dst, status string) (ordersResponse, error) {
	u := c.baseUrl + fmt.Sprintf(
		`/market/orders/list?srcCurrency=%s&dstCurrency=%s&details=2&status=%s&type=%s`,
		src, dst, status, side)

	var v ordersResponse
	if err := c.do(http.MethodGet, u, c.authHeader(), nil, &v); err != nil {
		return ordersResponse{}, err
	}
	return v, nil
}

func (c *httpClient) updateOrder(id, status string) (updateOrderResponse, error) {
	u := c.baseUrl + "/market/orders/update-status"
	body, _ := json.Marshal(map[string]interface{}{
		"id":     id,
		"status": status,
	})
	var v updateOrderResponse
	err := c.do(http.MethodPost, u, c.authHeader(), bytes.NewReader(body), &v)
	if err != nil {
		return updateOrderResponse{}, err
	}
	return v, nil
}

func (c *httpClient) cancelOrder(order order.Order) (cancelOrderResponse, error) {
	u := fmt.Sprintf("%s/positions/%s/close", c.baseUrl, strconv.Itoa(int(order.Id)))
	body, _ := json.Marshal(map[string]interface{}{
		"amount": order.Amount.String(),
		"price":  order.Price.String(),
	})
	var v cancelOrderResponse
	err := c.do(http.MethodPost, u, c.authHeader(), bytes.NewReader(body), &v)
	if err != nil {
		return cancelOrderResponse{}, err
	}
	return v, nil
}

func (c *httpClient) trades(src, dst string, fromId int) (tradesResponse, error) {
	u := c.baseUrl + fmt.Sprintf(
		`/market/trades/list?srcCurrency=%s&dstCurrency=%s&fromId=%d`,
		src, dst, fromId)

	if fromId == 0 {
		u = c.baseUrl + fmt.Sprintf(
			`/market/trades/list?srcCurrency=%s&dstCurrency=%s`,
			src, dst)
	}

	var v tradesResponse
	if err := c.do(http.MethodGet, u, c.authHeader(), nil, &v); err != nil {
		return tradesResponse{}, err
	}
	return v, nil
}

func (c *httpClient) do(method, url string, header http.Header, body io.Reader, v interface{}) error {
	req, err := http.NewRequest(method, url, body)
	req.Header = header
	if err != nil {
		return err
	}
	res, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *httpClient) authHeader() http.Header {
	header := http.Header{}
	header.Set("Authorization", "Token "+c.apiKey)
	return header
}
