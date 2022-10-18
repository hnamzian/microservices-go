package data

import (
	"testing"
)

func TestProductValidation(t *testing.T) {
	prod := &Product{
		Name: "tea",
		Description: "Nice tea",
		Price: 1.56,
		SKU: "abc-def-ghi",
	}

	err := prod.ValidateProduct()
	if err != nil {
		t.Fatal(err)
	}
}
