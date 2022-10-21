package server

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	currency "github.com/hnamzian/microservices-go/currency/currency"
)

type Currency struct {
	log hclog.Logger
	currency.UnimplementedCurrencyServer
}

func NewCurrencyServer(log hclog.Logger) *Currency {
	return &Currency{
		log: log,
		UnimplementedCurrencyServer: currency.UnimplementedCurrencyServer{},
	}
}

func (cs *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	cs.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	return &currency.RateResponse{Rate: 0.5}, nil
}
