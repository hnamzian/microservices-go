// Package classifications Product API
//
// Documentation for Product API
//
//	 Schemes: http
//		BasePath: /
//		Version: 1.0.0
//
//		Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//
// swagger:meta
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-hclog"
	currency "github.com/hnamzian/microservices-go/product-api/currency"
	"github.com/hnamzian/microservices-go/product-api/handlers"
	"google.golang.org/grpc"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "product-api",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// create grpc connection
	gconn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		l.Error("Unable to creat grpc client", "error", err)
	}

	defer gconn.Close()

	// create currency client
	cc := currency.NewCurrencyClient(gconn)


	ph := handlers.NewProducts(l, cc)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetOne)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.Create)
	postRouter.Use(ph.Middleware)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.Update)
	putRouter.Use(ph.Middleware)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	ch := gohandlers.CORS(
		gohandlers.AllowedOrigins([]string{"http://localhost:3000"}),
	)

	s := &http.Server{
		Addr:    ":9090",
		Handler: ch(sm),

		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Start Running Server on Port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to Start Server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Info("Received terminate, graceful shutdown: %s", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
