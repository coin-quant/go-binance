package binance

import (
	"context"
	"fmt"
	"net/http"
)

type MarginDataEventType string

const (
	baseMarginWsMainURL = "wss://margin-stream.binance.com/ws"

	MarginDataEventTypeUserLiabilityChange     MarginDataEventType = "USER_LIABILITY_CHANGE"
	MarginDataEventTypeMarginLevelStatusChange MarginDataEventType = "MARGIN_LEVEL_STATUS_CHANGE"

	SideEffectTypeAutoBorrowRepay SideEffectType = "AUTO_BORROW_REPAY"
)

// StartUserStreamService create listen key for user stream service
type StartMarginAccountStreamService struct {
	c *Client
}

// Do send request
func (s *StartMarginAccountStreamService) Do(ctx context.Context, opts ...RequestOption) (listenKey string, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/margin/listen-key",
		secType:  secTypeAPIKey,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	j, err := newJSON(data)
	if err != nil {
		return "", err
	}
	listenKey = j.Get("listenKey").MustString()
	return listenKey, nil
}

// KeepaliveUserStreamService update listen key
type KeepaliveMarginAccountStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveMarginAccountStreamService) ListenKey(listenKey string) *KeepaliveMarginAccountStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveMarginAccountStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: "/sapi/v1/margin/listen-key",
		secType:  secTypeAPIKey,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService delete listen key
type CloseMarginAccountStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *CloseMarginAccountStreamService) ListenKey(listenKey string) *CloseMarginAccountStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseMarginAccountStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/sapi/v1/margin/listen-key",
		secType:  secTypeAPIKey,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
}

type GetAvailableInventoryService struct {
	c          *Client
	marginType string
}

type AvailableInventory struct {
	Assets     map[string]string `json:"assets"`
	UpdateTime int64             `json:"updateTime"`
}

func (s *GetAvailableInventoryService) MarginType(marginType string) *GetAvailableInventoryService {
	s.marginType = marginType
	return s
}

// Do send request
func (s *GetAvailableInventoryService) Do(ctx context.Context, opts ...RequestOption) (res *AvailableInventory, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/margin/available-inventory",
		secType:  secTypeSigned,
	}
	r.setParam("type", s.marginType)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(AvailableInventory)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NewStartUserStreamService init starting user stream service
func (c *Client) NewStartMarginAccountStreamService() *StartMarginAccountStreamService {
	return &StartMarginAccountStreamService{c: c}
}

// NewKeepaliveUserStreamService init keep alive user stream service
func (c *Client) NewKeepaliveMarginAccountStreamService() *KeepaliveMarginAccountStreamService {
	return &KeepaliveMarginAccountStreamService{c: c}
}

// NewCloseUserStreamService init closing user stream service
func (c *Client) NewCloseMarginAccountStreamService() *CloseMarginAccountStreamService {
	return &CloseMarginAccountStreamService{c: c}
}

func (c *Client) NewGetAvailableInventoryService() *GetAvailableInventoryService {
	return &GetAvailableInventoryService{c: c}
}

type WsMarginDataHandler func(event *WsMarginDataEvent)

// WsUserDataEvent define user data event
type WsMarginDataEvent struct {
	Event     MarginDataEventType `json:"e"`
	Time      int64               `json:"E"`
	Asset     string              `json:"a"`
	Type      string              `json:"t"`
	Principal string              `json:"p"`
	Interest  string              `json:"i"`
}

// WsMarginUserDataServe serve user data handler with listen key
func WsMarginDataServe(listenKey string, handler WsMarginDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", baseMarginWsMainURL, listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarginDataEvent)
		err = json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
