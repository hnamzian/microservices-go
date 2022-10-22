package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-api/data"
)

type Products struct {
	log   hclog.Logger
	pdb *data.ProductsDB
}

func NewProducts(l hclog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, pdb}
}

type KeyProduct struct{}

func (p *Products) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.log.Info("Middleware")

		prod := &data.Product{}
		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.log.Error("Could not parse product from body", "error", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "Could not parse product from body"}, rw)
			return
		}

		err = prod.ValidateProduct()
		if err != nil {
			p.log.Error("Unable to parse Product", "error", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: fmt.Sprintf("Unable to parse Product: %s", err)}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		h.ServeHTTP(rw, req)
	})
}

type GenericError struct {
	Message string `json:"message"`
}