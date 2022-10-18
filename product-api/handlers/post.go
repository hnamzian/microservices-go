package handlers

import (
	"microservices-go/data"
	"net/http"
)

func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("[DEBUG] Create New Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}