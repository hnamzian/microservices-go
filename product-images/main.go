package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-images/configs"
	"github.com/hnamzian/microservices-go/product-images/files"
	"github.com/hnamzian/microservices-go/product-images/handlers"
)

func main() {
	cfg := configs.Load()

	hcl := hclog.New(&hclog.LoggerOptions{
		Name:  "product-images",
		Level: hclog.Debug,
	})

	fl, err := files.NewLocal(hcl, cfg.BasePath)
	if err != nil {
		hcl.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	fh := handlers.NewFiles(hcl, fl)
	sm := mux.NewRouter()

	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{3}}", fh.SaveFile)

	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{3}}",
		http.StripPrefix("/images", http.FileServer(http.Dir(cfg.BasePath))))

	hcl.Info("Server Starts Running", "bindAddress", cfg.BindAddress)

	s := &http.Server{
		Addr:         cfg.BindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()

}
