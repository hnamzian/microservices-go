package handlers

import (
	"context"
	"log"
	"microservices-go/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("Handle Get Products")

	lp := data.GetProductList()

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Could not Marshal product list", http.StatusInternalServerError)
	// }
	// rw.Write([]byte(d))

	// 2nd method: use json Encoder module which is faster and does not nedd any buffer ot local vars
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Could not encode product list into json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("Handle PUT Products")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(rw, "Invalid Product Id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
}

type KeyProduct struct {}

func (p *Products) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Printf("Middleware")

		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Could not parse product from body", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		h.ServeHTTP(rw, req)
	})
}