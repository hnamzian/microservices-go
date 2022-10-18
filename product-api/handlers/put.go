package handlers

import (
	"microservices-go/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("[DEBUG] Update A Product")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(rw, "[ERROR] Invalid Product Id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "[ERROR] Product Not Found", http.StatusNotFound)
		return
	}
}