package data

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/hashicorp/go-hclog"
	"github.com/hnamzian/microservices-go/product-api/currency"
)

var ErrorProductNotFound = fmt.Errorf("Product Not Found")

type ProductsDB struct {
	log hclog.Logger
	cc  currency.CurrencyClient
}

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


func NewProductsDB(log hclog.Logger, cc currency.CurrencyClient) *ProductsDB {
	return &ProductsDB{ log, cc	}
}

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

func (pd *ProductsDB) GetProductById(id int, dest string) (*Product, error) {
	pid := findIndexProductById(id)
	if pid == -1 {
		return nil, ErrorProductNotFound
	}

	rate, err := pd.getRate(dest)
	if err != nil {
		pd.log.Error("Unable to get rate", "currency", dest, "error", err)
		return nil, err
	}

	np := *productList[pid]
	np.Price = np.Price * rate

	return &np, nil
}

func (pd *ProductsDB) GetProductsAll(dest string) (Products, error) {
	pd.log.Debug("dest", "currency", dest)
	products := Products{}
	for _, p := range productList {
		np := *p

		rate, err := pd.getRate(dest)
		if err != nil {
			return nil, err
		}

		np.Price = np.Price * rate

		products = append(products, &np)
	}

	return products, nil
}

func (pd *ProductsDB) AddProduct(p *Product) {
	p.ID = getNextId()
	// p.CreatedOn = time.Now().UTC().String()
	// p.UpdatedOn = time.Now().UTC().String()
	productList = append(productList, p)
}

func (pd *ProductsDB) UpdateProduct(id int, p *Product) error {
	pid := findIndexProductById(id)

	if pid == -1 {
		return ErrorProductNotFound
	}

	p.ID = id
	productList[pid] = p
	return nil
}

func (pd *ProductsDB) DeleteProduct(id int) error {
	pid := findIndexProductById(id)
	if pid == -1 {
		return ErrorProductNotFound
	}

	productList = append(productList[:pid], productList[pid+1])

	return nil
}

func getNextId() int {
	id := productList[len(productList)-1].ID
	return id + 1
}

func findIndexProductById(id int) (int) {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (pd *ProductsDB) getRate(dest string) (float64, error) {
	rr := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["EUR"]),
		Destination: currency.Currencies(currency.Currencies_value[dest]),
	}
	rate, err := pd.cc.GetRate(context.Background(), rr)

	pd.log.Info("rate", "dest", dest, "rate", rate.Rate)

	return rate.Rate, err
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
