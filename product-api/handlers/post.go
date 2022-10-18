package handlers

import (
	"microservices-go/data"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse
// Create handles POST requests to add new products
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("[DEBUG] Create New Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}