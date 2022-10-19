package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-images/configs"
	"github.com/hnamzian/microservices-go/product-images/files"
	"github.com/hnamzian/microservices-go/product-images/handlers"
)

func main() {
	cfg := configs.Load()

	l := log.New(os.Stdout, "product-images", log.LstdFlags)

	fl, _ := files.NewLocal(l, cfg.BasePath)

	fh := handlers.NewFiles(l, fl)
	sm := mux.NewRouter()

	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{3}}", fh.SaveFile)

	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-zA-Z]{3}}",
		http.StripPrefix("/images", http.FileServer(http.Dir(cfg.BasePath))))

	s := &http.Server{
		Addr:         cfg.BindAddress,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()
}
