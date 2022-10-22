package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-api/data"
)

// swagger:route GET /products products listProducts
//
// # Lists products in the system
//
// This will show all available products by default.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http
//
//	Responses:
//	  200: productsResponse
//	  422: validationError
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("Get All Products")

	params := mux.Vars(r)
	dest := params["currency"]

	lp, err := p.pdb.GetProductsAll(dest)
	if err != nil {
		p.log.Error("Could not get products", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Could not get products"}, rw)
		return
	}

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Could not Marshal product list", http.StatusInternalServerError)
	// }
	// rw.Write([]byte(d))

	// 2nd method: use json Encoder module which is faster and does not nedd any buffer ot local vars
	err = data.ToJSON(lp, rw)
	if err != nil {
		p.log.Error("Could not encode product list into json", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Could not encode product list into json"}, rw)
		return
	}
}

func (p *Products) GetOne(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("Get One product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		p.log.Error("Unable to parse Id", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Unable to parse Id"}, rw)
		return
	}

	dest := r.URL.Query().Get("currency")

	product, err := p.pdb.GetProductById(id, dest)

	switch err {
	case nil:

	case data.ErrorProductNotFound:
		p.log.Error("Product Not Found", "error", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product Not Found"}, rw)
		return

	default:
		p.log.Error("Unable to fetching product", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Unable to fetching product"}, rw)
		return
	}

	err = data.ToJSON(product, rw)
	if err != nil {
		p.log.Error("Could not encode product list into json", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Could not encode product list into json"}, rw)
		return
	}
}
