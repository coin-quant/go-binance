package delivery

import (
	"context"
	"encoding/json"
	"github.com/adshao/go-binance/v2/common"
	"net/http"
)

// PremiumIndexService get premium index
type PremiumIndexService struct {
	c      *Client
	symbol *string
	pair   *string
}

// PremiumIndex define premium index of mark price
type PremiumIndex struct {
	Symbol               string `json:"symbol"`
	Pair                 string `json:"pair"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	InterestRate         string `json:"interestRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	Time                 int64  `json:"time"`
}

// Symbol set symbol
func (s *PremiumIndexService) Symbol(symbol string) *PremiumIndexService {
	s.symbol = &symbol
	return s
}

func (s *PremiumIndexService) Pair(pair string) *PremiumIndexService {
	s.pair = &pair
	return s
}

// Do send request
func (s *PremiumIndexService) Do(ctx context.Context, opts ...RequestOption) (res []*PremiumIndex, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/premiumIndex",
		secType:  secTypeNone,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.pair != nil {
		r.setParam("pair", *s.pair)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	data = common.ToJSONList(data)
	if err != nil {
		return []*PremiumIndex{}, err
	}
	res = make([]*PremiumIndex, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*PremiumIndex{}, err
	}
	return res, nil
}

func (c *Client) NewPremiumIndexService() *PremiumIndexService {
	return &PremiumIndexService{c: c}
}
