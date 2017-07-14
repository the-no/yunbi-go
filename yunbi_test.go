package yunbi

import (
	"fmt"
	"testing"
)

func TestDepth(t *testing.T) {
	m, err := Depth(&DepthRequest{"btscny", 0})

	var asknum float64
	var tprice float64
	for _, ask := range m.Asks {
		p, _ := ask[0].Float64()
		n, _ := ask[1].Float64()
		asknum += n
		tprice += n * p

	}
	fmt.Println(asknum, tprice, tprice/asknum)

	var bnum float64
	var bprice float64
	for _, ask := range m.Bids {
		p, _ := ask[0].Float64()
		n, _ := ask[1].Float64()
		bnum += n
		bprice += n * p

	}
	fmt.Println(bnum, bprice, bprice/bnum)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(m)
	}
}

/*

func TestGetTicker(t *testing.T) {
	m, err := MarketTicker("btccny")
	fmt.Println(err)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(m)
	}
}

func TestTickerList(t *testing.T) {
	markets, err := TickerList()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(markets)
	}
}

func TestOrderBookList(t *testing.T) {
	resp, err := OrderBookList(&OrderBookRequest{"btccny", 10, 10})
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(resp)
	}
}

func TestDepth(t *testing.T) {
	resp, err := Depth(&DepthRequest{"btccny", 10})
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(resp)
	}
}

func TestServerTimestamp(t *testing.T) {
	resp, err := ServerTimestamp()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(resp)
	}
}

func TestAccount(t *testing.T) {
	resp, err := AccountInfo()
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(resp)
	}
}
*/

/*func TestDepositList(t *testing.T) {
	err := DepositList(&DepositsRequest{})
	if err != nil {
		t.Error(err.Error())
	} else {
		//t.Log(resp)
	}
}
*/

/*func TestCreateOrder(t *testing.T) {
	req := &OrderRequest{Side: "sell", Market: "btsscny", Volume: 18.5814, Price: 2.1589}
	resp, err := CreateOrder(req)
	fmt.Println(resp, err)
}
*/
/*func TestClearOrder(t *testing.T) {

	resp, err := ClearOrder("buy")
	fmt.Println(resp, err)
}*/
