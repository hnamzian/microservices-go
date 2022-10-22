package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

var ErrorProductNotFound = fmt.Errorf("Product Not Found")

// Product defines the structure for an API product
// swagger:model Product
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}


type Products []*Product

func (p *Product) ValidateProduct() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func GetOneProduct(id int) (*Product, error) {
	_, pos, err := findProduct(id)
	if err != nil {
		return nil, err
	}

	np := *productList[pos]

	return &np, nil
}

func GetProductList() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	// p.CreatedOn = time.Now().UTC().String()
	// p.UpdatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

func DeleteProduct(id int) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return ErrorProductNotFound
	}

	productList = append(productList[:pos], productList[pos+1])

	return nil
}

func getNextId() int {
	id := productList[len(productList)-1].ID
	return id + 1
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrorProductNotFound
}

var productList = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		// CreatedOn:   time.Now().UTC().String(),
		// UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		// CreatedOn:   time.Now().UTC().String(),
		// UpdatedOn:   time.Now().UTC().String(),
	},
}
