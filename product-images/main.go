package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-images/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-images", log.LstdFlags)

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
		l.Println("Server starts Running on Port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	signal := <- sigChan
	l.Printf("Received terminate, graceful shutdown:%s", signal)

	tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(tc)
}
