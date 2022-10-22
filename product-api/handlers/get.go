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
	p.l.Info("[DEBUG] Get All Products")

	params := mux.Vars(r)
	dest := params["currency"]

	lp, err := p.pdb.GetProductsAll(dest)
	if err != nil {
		http.Error(rw, "[ERROR] Could not get products", http.StatusInternalServerError)
		return
	}

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Could not Marshal product list", http.StatusInternalServerError)
	// }
	// rw.Write([]byte(d))

	// 2nd method: use json Encoder module which is faster and does not nedd any buffer ot local vars
	err = lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "[ERROR] Could not encode product list into json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) GetOne(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("[DEBUG] Get One product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		p.l.Error("[ERROR] Unable to parse Id", "error", err)
		http.Error(rw, "Unable to parse Id", http.StatusInternalServerError)
		return
	}

	dest := r.URL.Query().Get("currency")

	product, err := p.pdb.GetProductById(id, dest)

	switch err {
	case nil:

	case data.ErrorProductNotFound:
		p.l.Error("[ERROR] Product Not Found", "error", err)
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return

	default:
		p.l.Error("Unable to fetching product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		http.Error(rw, "Unable to fetching product", http.StatusNotFound)
		return
	}

	err = product.ToJSON(rw)
	if err != nil {
		http.Error(rw, "[ERROR] Could not encode product list into json", http.StatusInternalServerError)
		return
	}
}
