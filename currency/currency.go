package currency

import (
	"encoding/csv"
	"strconv"
	"strings"

	"net/http"
)

type Currencies struct {
	currenciesList map[string]float64
}

func New() (currencies *Currencies) {
	return new(Currencies)
}

func Load(currs ...interface{}) (currencies *Currencies, err error) {
	currencies = new(Currencies)
	currencies.currenciesList = map[string]float64{}
	cs := []string{}

	for _, c := range currs {
		cs = append(cs, c.(string))
	}
	m, err := retreiveCurrencies(cs)
	if err != nil {
		return
	}
	currencies.currenciesList = m
	return
}

func (c *Currencies) Get(currency string) (rate float64, err error) {

	rate, ok := c.currenciesList[currency]
	if !ok {
		var cmap map[string]float64
		cmap, err = retreiveCurrencies([]string{currency})
		rate = cmap[currency]
		c.currenciesList[currency] = rate
	}
	return rate, err
}

func retreiveCurrencies(currList []string) (ratemap map[string]float64, err error) {
	ratemap = map[string]float64{}
	response, err := http.Get("http://download.finance.yahoo.com/d/quotes.csv?s=" + strings.Join(currList, "=X,") + "=X&f=c4l1&e=.csv")
	if nil != err {
		return
	}

	enc := csv.NewReader(response.Body)
	lines, err := enc.ReadAll()
	if nil != err {
		return
	}

	for _, line := range lines {
		var rate float64
		rate, err = strconv.ParseFloat(line[1], 64)
		if nil != err {
			return
		}
		ratemap[line[0]] = rate
	}
	return
}
