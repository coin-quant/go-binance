package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type WsCombinedBookTickerEvent struct {
	Data   *WsBookTickerEvent `json:"data"`
	Stream string             `json:"stream"`
}

func WsCombinedBookTickerServe(symbols []string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@bookTicker", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsCombinedBookTickerEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event.Data)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

type WsRawBookTickerHandler func(raw []byte)

func WsRawBookTickerServe(symbol string, handler WsRawBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@bookTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		handler(message)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func WsCombinedRawBookTickerServe(symbols []string, handler WsRawBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@bookTicker", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		handler(message)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

type WsCombinedMarkPriceEvent struct {
	Data   *WsMarkPriceEvent `json:"data"`
	Stream string            `json:"stream"`
}

type WsRawMarkPriceHandler func([]byte)

func wsCombinedRawMarkPriceServe(endpoint string, handler WsRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		handler(message)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func WsCombinedRawMarkPriceServe(symbols []string, handler WsRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@markPrice", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	return wsCombinedRawMarkPriceServe(endpoint, handler, errHandler)
}

func WsCombinedRawMarkPriceServeWithRate(symbolLevels map[string]time.Duration, handler WsRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, rate := range symbolLevels {
		var rateStr string
		switch rate {
		case 3 * time.Second:
			rateStr = ""
		case 1 * time.Second:
			rateStr = "@1s"
		default:
			return nil, nil, fmt.Errorf("invalid rate. Symbol %s (rate %d)", symbol, rate)
		}

		endpoint += fmt.Sprintf("%s@markPrice%s", strings.ToLower(symbol), rateStr) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	return wsCombinedRawMarkPriceServe(endpoint, handler, errHandler)
}

type WsAllRawMarkPriceHandler func([]byte)

func wsAllRawMarkPriceServe(endpoint string, handler WsAllRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		handler(message)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func WsAllRawMarkPriceServe(handler WsAllRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!markPrice@arr", getWsEndpoint())
	return wsAllRawMarkPriceServe(endpoint, handler, errHandler)
}

func WsAllRawMarkPriceServeWithRate(rate time.Duration, handler WsAllRawMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, nil, errors.New("Invalid rate")
	}
	endpoint := fmt.Sprintf("%s/!markPrice@arr%s", getWsEndpoint(), rateStr)
	return wsAllRawMarkPriceServe(endpoint, handler, errHandler)
}
