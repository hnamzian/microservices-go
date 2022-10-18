package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hnamzian/microservices-go/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func (p *Products) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Printf("Middleware")

		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Could not parse product from body", http.StatusBadRequest)
			return
		}

		err = prod.ValidateProduct()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Unable to parse Product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		h.ServeHTTP(rw, req)
	})
}
