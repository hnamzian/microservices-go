package server

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	currency "github.com/hnamzian/microservices-go/currency/currency"
	data "github.com/hnamzian/microservices-go/currency/data"
)

type Currency struct {
	log hclog.Logger
	rates *data.ExchangeRates

	currency.UnimplementedCurrencyServer
}

func NewCurrencyServer(log hclog.Logger, rates *data.ExchangeRates) *Currency {
	return &Currency{
		log: log,
		rates: rates,
		UnimplementedCurrencyServer: currency.UnimplementedCurrencyServer{},
	}
}

func (cs *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	cs.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := cs.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &currency.RateResponse{Rate: rate}, nil
}
