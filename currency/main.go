package main

import (
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/currency/currency"
	"github.com/hnamzian/microservices-go/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	

	gs := grpc.NewServer()

	cs := server.NewCurrencyServer(log)

	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to create listener", "error", err)
	}
	gs.Serve(listener)
}
