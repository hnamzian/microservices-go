package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewExchangeRates(log hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{
		log:   log,
		rates: map[string]float64{},
	}

	err := er.getRates()
	
	return er, err
}

func (er *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to get rates")
	}

	er.log.Info("response", "status", resp.StatusCode)

	md := &Cubes{}
	err = xml.NewDecoder(resp.Body).Decode(md)
	if err != nil {
		return err
	}

	for _, c := range md.CubeData {
		rf, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		er.rates[c.Currency] = rf
	}

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
