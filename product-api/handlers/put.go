package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hnamzian/microservices-go/product-api/data"
)

func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("Update A Product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "[ERROR] Invalid Product Id"}, rw)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = p.pdb.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "[ERROR] Product Not Found"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
