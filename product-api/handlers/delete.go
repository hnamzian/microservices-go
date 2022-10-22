package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
//
// # Delete a product from datastore
//
// This will delete a product from datastore
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
//	  201: noContent
//	  422: validationError
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("[DEBUG] Delete a product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		p.l.Info("[ERROR] Invalid Product Id")
		http.Error(rw, "Invalid product Id", http.StatusBadRequest)
		return
	}

	err = p.pdb.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		p.l.Info("[ERROR] Deleting product id does not exist")

		http.Error(rw, fmt.Sprintf("%s", err), http.StatusNotFound)
	}

	if err != nil {
		p.l.Info("[ERROR] Deleting product", err)

		http.Error(rw, fmt.Sprintf("Internal Error: %s", err), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
