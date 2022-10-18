// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/hnamzian/microservices-go/product-api/data"

// A list of products
// swagger:response productsResponse
type ProductsResponseWrapper struct {
	// The error message
	// in: body
	Body []data.Product
}
