package yunbi

import (
	"encoding/json"
	"time"
)

var (
	AccessKey = ""
	SecretKey = ""
)

const (
	YUN_BI_HOST = "https://yunbi.com"
)

type Ticker struct {
	Buy  string `json:"buy"`  // 当前买入价
	Sell string `json:"sell"` // 当前卖出价
	Low  string `json:"low"`  // 过去24小时之内的最低价
	High string `json:"high"` // 过去24小时之内的最高价
	Last string `json:"last"` // 最后成交价
	Vol  string `json:"vol"`  // 过去24小时之内的总成交量
}

type Market struct {
	At     int64  `json:"at"` // 以秒为单位的时间戳
	Ticker Ticker `json:"ticker"`
}

type Account struct {
	Currency string `json:"currency"` // 账户的币种 例如 cny 或 btc
	Balance  string `json:"balance"`  // 账户余额, 不包括冻结资金
	Locked   string `json:"locked"`   // 账户冻结资金
}

type Member struct {
	Sn        string     `json:"sn"`        // 用户的唯一编号
	Name      string     `json:"name"`      // 用户名字
	Email     string     `json:"email"`     // 用户邮件地址
	Activated bool       `json:"activated"` // 用户是否已激活
	Accounts  []*Account `json:"accounts"`  // 用户的所有账户信息 [参见 Account]
}

type Trade struct {
	Id        int64     `json:"id"`         // 交易的唯一ID
	Price     string    `json:"price"`      // 成交价
	Volume    string    `json:"volume"`     // 成交数量
	Market    string    `json:"market"`     // 交易所属的市场
	CreatedAt time.Time `json:"created_at"` // 成交时间
	OrderId   int64     `json:"order_id"`   // 交易所属 Order ID
	Side      string    `json:"side"`       // buy/sell 买或者卖
}

type Order struct {
	Id       int64  `json:"id"`        // 唯一的 Order ID
	Side     string `json:"side"`      // Buy/Sell 代表买单/卖单
	Price    string `json:"price"`     // 出价
	AvgPrice string `json:"avg_price"` // 平均成交价
	State    string `json:"state"`     // 订单的当前状态 [wait,done,cancel]  wait  表明订单正在市场上挂单
	// 是一个active order 此时订单可能部分成交或者尚未成交 done   代表订单已经完全成交 cancel 代表订单已经被撤销
	Market          string    `json:"market"`           // 订单参与的交易市场
	CreatedAt       time.Time `json:"created_at"`       // 下单时间 ISO8601格式
	Volume          string    `json:"volume"`           // 购买/卖出数量
	RemainingVolume string    `json:"remaining_volume"` // 还未成交的数量 remaining_volume 总是小于等于 volume 在订单完全成交时变成 0
	ExecutedVolume  string    `json:"executed_volume"`  // 已成交的数量  volume = remaining_volume + executed_volume
	TradesCount     int       `json:"trades_count"`     // 订单的成交数 整数值 未成交的订单为 0 有一笔成交的订单为 1 通过该字段可以判断订单是否处于部分成交状态
	Trades          []*Trade  `json:"trades"`           // 订单的详细成交记录 参见Trade 注意: 只有某些返回详细订单数据的 API 才会包含 Trade 数据
}

type Voucher struct {
	OrderId int64  `json:"order_id"` // 订单 id
	AskFee  string `json:"ask_fee"`  // 卖单交易费用，如果该挂单是买单，则卖单交易费为0
	BidFee  string `json:"bid_fee"`  // 买单交易费用，如果该挂单是卖单，则买单交易费为0
}

type OrderBook struct {
	Asks []int64 `json:"asks"` // 卖单列表
	Bids []int64 `json:"bids"` // 买单列表
}

type OrderBookRequest struct {
	Market    string
	AsksLimit int
	BidsLimit int
}

type OrderBookResponse struct {
	Asks []*Order `json:"asks"`
	Bids []*Order `json:"bids"`
}

type DepthRequest struct {
	Market string
	Limit  int
}

type DepthResponse struct {
	Timestamp int64           `json:"timestamp"`
	Asks      [][]json.Number `json:"asks"`
	Bids      [][]json.Number `json:"bids"`
}

type TradeRequest struct {
	Market    string
	Limit     int
	Timestamp time.Time
	From      time.Time
	To        time.Time
	OrderBy   string
}

type TradeResponse struct {
	Trades []*Trade `json:"trades"`
}

type KlineRequest struct {
	Market    string
	Limit     int
	Timestamp time.Time
	Period    int
}

type KlineResponse struct {
	Lines [][]int64 `json:"lines"`
}

type AddressResponse struct {
	Messages string `json:"messages"`
	Ip       string `json:"ip"`
}

type KWithTradesRequest struct {
	Market    string
	Limit     int
	Timestamp time.Time
	Period    int
	TradeId   int64
}

type KWithTradesResponse struct {
	K      [][]int64 `json:"k"`
	Trades []*Trade  `json:"trades"`
	//Err    YunBiError `json:"error"`
}

type MarketInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type YunBiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (self YunBiError) Error() string {
	return self.Error()
}

type DepositsRequest struct {
	Currency string
	Limit    int
	State    string
}

type OrderRequest struct {
	Side      string  // Buy/Sell 代表买单/卖单
	Price     float64 // 出价
	Market    string  // 订单参与的交易市场
	Volume    float64 // 购买/卖出数量
	OrderType string
}

type OrderResponse struct {
	Id int64 `json:"id"` // Buy/Sell 代表买单/卖单
}

type ClearOrderResponse struct {
	Orders Order `json:"orders"` // Buy/Sell 代表买单/卖单
}
