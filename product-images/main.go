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
	"github.com/hnamzian/microservices-go/product-images/files"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	bindAddress := os.Getenv("BIND_ADDRESS")
	basePath := os.Getenv("BASE_PATH")

	l := hclog.New(&hclog.LoggerOptions{
		Name:  "product-images",
		Level: hclog.LevelFromString("DEBUG"),
	})

	stor, err := files.NewLocal(basePath, 1024*5*1000)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}


	ph := handlers.NewFiles(stor, l)

	sm := mux.NewRouter()
	sm.Methods(http.MethodPost).Subrouter().HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", ph.ServeHTTP)
	sm.Methods(http.MethodGet).Subrouter().Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", 
		http.StripPrefix("/images", http.FileServer(http.Dir(basePath))))


	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := http.Server{
		Addr:    bindAddress,
		Handler: ch(sm),

		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Server starts Running", "bindAddress", bindAddress)
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
