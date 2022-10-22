package data

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func Test_NewRates(t *testing.T) {
	tr, err := NewExchangeRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Rates %#v", tr.rates)
}
