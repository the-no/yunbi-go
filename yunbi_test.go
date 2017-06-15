package yunbi

import (
	"testing"
)

func TestGetTicker(t *testing.T) {
	m, err := MarketTicker("btccny")
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


func TestDepositList(t *testing.T) {
	err := DepositList(&DepositsRequest{})
	if err != nil {
		t.Error(err.Error())
	} else {
		//t.Log(resp)
	}
}
