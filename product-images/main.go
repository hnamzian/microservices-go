package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-images/handlers"
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "product-images",
		Level: hclog.LevelFromString("DEBUG"),
	})

	handlers.NewFiles(l)

	fm := mux.NewRouter()

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := http.Server{
		Addr:    ":9090",
		Handler: ch(fm),

		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Server starts Running on Port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to Start Server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	signal := <-sigChan
	l.Info("Received terminate, graceful shutdown:%s", signal)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
