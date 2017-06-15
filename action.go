package yunbi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func YunBiSignature(mothod, uri, query, key string) string {
	params := []string{mothod, uri}
	if query != "" {
		params = append(params, query)
	}

	payload := strings.Join(params, "|")
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	data := mac.Sum(nil)
	return hex.EncodeToString(data)
}

func createRequest(mothod, uri string, params map[string]string) *http.Request {

	request, err := http.NewRequest(mothod, YUN_BI_HOST+uri, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Content-Type", "charset=utf-8")
	request.Close = true

	q := url.Values{}
	if params != nil {
		for k, v := range params {
			q.Add(k, v)
		}
		if params["access_key"] != "" {
			tonce := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
			q.Add("tonce", tonce)
			query := q.Encode()
			sign := YunBiSignature(mothod, uri, query, SecretKey)
			q.Add("signature", sign)
			//params["tonce"] = tonce
			//params["signature"] = sign
		}
	}
	if mothod == "GET" {
		request.URL.RawQuery = q.Encode()
	} else {
		/*data, err := json.Marshal(params)
		if err != nil {
			return nil
		}
		request.Body = ioutil.NopCloser(bytes.NewReader(data))*/
		request.PostForm = q
	}
	return request
}

func getHttpsClient() (client *http.Client) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	return
}

func MarketList() (markets []*MarketInfo, err error) {
	req := createRequest("GET", "/api/v2/markets.json", nil)
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	markets = []*MarketInfo{}
	err = json.NewDecoder(resp.Body).Decode(markets)
	if err != nil {
		return nil, err
	}
	return markets, err
}

func MarketTicker(m string) (market *Market, err error) {
	req := createRequest("GET", "/api/v2/tickers/"+m+".json", nil)
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	market = &Market{}
	err = json.NewDecoder(resp.Body).Decode(market)
	if err != nil {
		return nil, err
	}
	return market, err
}

func TickerList() (markets map[string]*Market, err error) {
	req := createRequest("GET", "/api/v2/tickers.json", nil)
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	markets = make(map[string]*Market)
	if err = json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		return nil, err
	}
	return markets, nil
}

func OrderBookList(req *OrderBookRequest) (response *OrderBookResponse, err error) {
	if req.AsksLimit == 0 {
		req.AsksLimit = 10
	}
	if req.BidsLimit == 0 {
		req.BidsLimit = 10
	}

	params := map[string]string{
		"market":     req.Market,
		"asks_limit": strconv.Itoa(req.AsksLimit),
		"bids_limit": strconv.Itoa(req.BidsLimit),
	}
	request := createRequest("GET", "/api/v2/order_book.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &OrderBookResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, err
}

func Depth(req *DepthRequest) (response *DepthResponse, err error) {

	if req.Limit == 0 {
		req.Limit = 10
	}

	params := map[string]string{
		"market": req.Market,
		"limit":  strconv.Itoa(req.Limit),
	}
	request := createRequest("GET", "/api/v2/depth.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &DepthResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, err
}

func ServerTimestamp() (timestamp int64, err error) {
	req := createRequest("GET", "/api/v2/timestamp.json", nil)
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}
	//market = &Market{}
	err = json.NewDecoder(resp.Body).Decode(&timestamp)
	if err != nil {
		return 0, err
	}
	fmt.Println(timestamp, err)
	return timestamp, err
}

func tradeList(req *TradeRequest, mode int) (response *TradeResponse, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	params := map[string]string{
		"market":   req.Market,
		"limit":    strconv.Itoa(req.Limit),
		"order_by": req.OrderBy,
	}
	if mode == 1 {
		params["access_key"] = AccessKey
	}
	if !req.Timestamp.IsZero() {
		params["timestamp"] = strconv.Itoa(int(req.Timestamp.Unix()))
	}
	if !req.From.IsZero() {
		params["timestamp"] = strconv.Itoa(int(req.From.Unix()))
	}
	if !req.To.IsZero() {
		params["timestamp"] = strconv.Itoa(int(req.To.Unix()))
	}
	request := createRequest("GET", "/api/v2/trades.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	trades := []*Trade{}
	if err = json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}
	response = &TradeResponse{Trades: trades}
	return response, err
}

func TradeList(req *TradeRequest) (response *TradeResponse, err error) {
	return tradeList(req, 0)
}

func MyTradeList(req *TradeRequest) (response *TradeResponse, err error) {
	return tradeList(req, 1)
}

func KlineList(req *KlineRequest) (response *KlineResponse, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Period == 0 {
		req.Period = 1
	}

	params := map[string]string{
		"market": req.Market,
		"limit":  strconv.Itoa(req.Limit),
		"period": strconv.Itoa(req.Period),
	}

	if !req.Timestamp.IsZero() {
		params["timestamp"] = strconv.Itoa(int(req.Timestamp.Unix()))
	}

	request := createRequest("GET", "/api/v2/k.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	lines := [][]int64{}
	if err = json.NewDecoder(resp.Body).Decode(&lines); err != nil {
		return nil, err
	}
	response = &KlineResponse{Lines: lines}
	return response, err
}

func Address(addr string) (response *AddressResponse, err error) {
	req := createRequest("GET", "/api/v2/addresses/"+addr+".json", nil)
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &AddressResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	return response, err
}

func KWithTradesList(req *KWithTradesRequest) (response *KWithTradesResponse, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Period == 0 {
		req.Period = 1
	}
	if req.TradeId < 1 {
		return nil, errors.New("trade_id Requested")
	}

	params := map[string]string{
		"market":   req.Market,
		"limit":    strconv.Itoa(req.Limit),
		"period":   strconv.Itoa(req.Period),
		"trade_id": strconv.Itoa(int(req.TradeId)),
	}

	if !req.Timestamp.IsZero() {
		params["timestamp"] = strconv.Itoa(int(req.Timestamp.Unix()))
	}

	request := createRequest("GET", "/v2/k_with_pending_trades.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	response = &KWithTradesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, err
}

func AccountInfo() (response *Account, err error) {
	req := createRequest("GET", "/api/v2/members/me.json", map[string]string{"access_key": AccessKey})
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &Account{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(d), err)
	return response, nil
}

func Deposit(txid int) (err error) {
	if txid == 0 {
		return errors.New("txid Requested")
	}
	req := createRequest("GET", "/api/v2/deposit.json", map[string]string{"access_key": AccessKey, "txid": strconv.Itoa(txid)})
	cli := getHttpsClient()
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	/*	response = &Account{}
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return nil, err
		}*/
	d, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(d), err)
	return nil
}

func DepositList(req *DepositsRequest) (err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	params := map[string]string{
		"limit":      strconv.Itoa(req.Limit),
		"access_key": AccessKey,
	}

	if req.Currency != "" {
		params["currency"] = req.Currency
	}
	if req.State != "" {
		params["state"] = req.State
	}
	requset := createRequest("GET", "/api/v2/deposits.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(requset)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	/*	response = &Account{}
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return nil, err
		}*/
	d, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(d), err, requset)
	return nil
}

func CreateOrder(req *OrderRequest) (response *OrderResponse, err error) {
	if req.Side != "buy" && req.Side != "sell" {
		return nil, errors.New("side Requested [buy ,sell]")
	}
	if req.Volume == 0 || req.Price == 0 {
		return nil, errors.New("[Volume ,Price] Requested ")
	}
	if req.Market == "" {
		return nil, errors.New("Market Requested")
	}

	params := map[string]string{
		"access_key": AccessKey,
		"market":     req.Market,
		"side":       req.Side,
		"volume":     strconv.FormatFloat(req.Volume, 'f', -1, 64),
		"price":      strconv.FormatFloat(req.Price, 'f', -1, 64),
	}
	if req.OrderType != "" {
		params["order_type"] = req.OrderType
	}
	request := createRequest("POST", "/api/v2/orders.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &OrderResponse{}
	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if checkError(doc) != nil {
		return nil, err
	}
	err = json.Unmarshal(doc, response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func checkError(doc []byte) error {
	if bytes.Contains(doc, []byte("error")) {
		err := &YunBiError{}
		json.Unmarshal(doc, err)
		return err
	}
	return nil
}

func ClearOrder(side string) (response *OrderResponse, err error) {
	if side != "buy" && side != "sell" {
		return nil, errors.New("side Requested [buy ,sell]")
	}

	params := map[string]string{
		"access_key": AccessKey,
		"side":       side,
	}

	request := createRequest("POST", "/api/v2/order/clear.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &OrderResponse{}
	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if checkError(doc) != nil {
		return nil, err
	}
	err = json.Unmarshal(doc, response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func DeleteOrder(id int64) (response *OrderResponse, err error) {

	params := map[string]string{
		"access_key": AccessKey,
		"id":         strconv.FormatInt(id, 10),
	}

	request := createRequest("POST", "/api/v2/order/delete.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &OrderResponse{}
	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if checkError(doc) != nil {
		return nil, err
	}
	err = json.Unmarshal(doc, response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func OrderInfo(id int64) (response *Order, err error) {

	params := map[string]string{
		"access_key": AccessKey,
		"id":         strconv.FormatInt(id, 10),
	}

	request := createRequest("GET", "/api/v2/order.json", params)
	cli := getHttpsClient()
	resp, err := cli.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	response = &Order{}
	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if checkError(doc) != nil {
		return nil, err
	}
	err = json.Unmarshal(doc, response)
	if err != nil {
		return nil, err
	}

	return response, err
}
