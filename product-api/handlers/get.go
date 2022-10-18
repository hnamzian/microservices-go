package handlers

import (
	"microservices-go/data"
	"net/http"
)

func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("[DEBUG] Get All Products")

	lp := data.GetProductList()

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Could not Marshal product list", http.StatusInternalServerError)
	// }
	// rw.Write([]byte(d))

	// 2nd method: use json Encoder module which is faster and does not nedd any buffer ot local vars
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "[ERROR] Could not encode product list into json", http.StatusInternalServerError)
	}
}