package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/currency/currency"
	"github.com/hnamzian/microservices-go/currency/data"
	"github.com/hnamzian/microservices-go/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	rates, err := data.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to Get Rates")
		os.Exit(1)
	}
	cs := server.NewCurrencyServer(log, rates)
	
	gs := grpc.NewServer()
	
	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to create listener", "error", err)
	}
	gs.Serve(listener)
}
