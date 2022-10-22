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
	p.log.Debug("Delete a product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		p.log.Error("Invalid Product Id", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "Invalid product Id"}, rw)
		return
	}

	err = p.pdb.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		p.log.Error("Deleting product id does not exist", "error", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("%s", err)}, rw)
		return
	}

	if err != nil {
		p.log.Error("Deleting product", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: fmt.Sprintf("Internal Error: %s", err)}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
