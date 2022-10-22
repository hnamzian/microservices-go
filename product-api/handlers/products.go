package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-api/data"
)

type Products struct {
	l   hclog.Logger
	pdb *data.ProductsDB
}

func NewProducts(l hclog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, pdb}
}

type KeyProduct struct{}

func (p *Products) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Info("Middleware")

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
