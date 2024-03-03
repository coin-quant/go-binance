package futures

import (
	"encoding/json"
	"fmt"
	"strings"
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
