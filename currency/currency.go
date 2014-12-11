package currency

import (
	"encoding/csv"
	"strconv"
	"strings"

	"net/http"
)

type Currencies struct {
	currencies_list map[string]float64
}

func New() (currencies *Currencies) {
	return new(Currencies)
}

func Load(currs ...interface{}) (currencies *Currencies, err error) {
	currencies = new(Currencies)
	currencies.currencies_list = map[string]float64{}
	cs := []string{}

	for _, c := range currs {
		cs = append(cs, c.(string))
	}
	m, err := retreive_currencies(cs)
	if err != nil {
		return
	}
	currencies.currencies_list = m
	return
}

func (c *Currencies) Get(currency string) (rate float64, err error) {

	rate, ok := c.currencies_list[currency]
	if !ok {
		var cmap map[string]float64
		cmap, err = retreive_currencies([]string{currency})
		rate = cmap[currency]
		c.currencies_list[currency] = rate
	}
	return rate, err
}

func retreive_currencies(curr_list []string) (ratemap map[string]float64, err error) {
	ratemap = map[string]float64{}
	response, err := http.Get("http://download.finance.yahoo.com/d/quotes.csv?s=" + strings.Join(curr_list, "=X,") + "=X&f=c4l1&e=.csv")
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
