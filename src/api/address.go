package api

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Address struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	City  string `json:"city"`
	State string `json:"state"`
	Zip   string `json:"zip"`
}

type AddressManager interface {
	Find(s string) []*Address
}

type SimpleAddressManager struct {
	index map[string][]*Address
}

func (m *SimpleAddressManager) Find(s string) []*Address {
	return m.index[s]
}

func NewSimpleAddressManager(adrs []*Address) *SimpleAddressManager {
	m := SimpleAddressManager{index: make(map[string][]*Address)}
	for _, adr := range adrs {
		for _, v := range getIndexValues(adr) {
			m.index[v] = append(m.index[v], adr)
		}
	}
	return &m
}

func getIndexValues(adr *Address) []string {
	var r []string
	lines := []string{
		adr.Line1,
		adr.City,
		adr.State,
		adr.Zip,
	}
	for _, l := range lines {
		for _, w := range getWords(l) {
			r = append(r, w)
		}
	}
	return r
}

func getWords(s string) []string {
	return strings.Split(s, " ")
}

func LoadFromFile(file string) ([]*Address, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var adrs []*Address
	err = json.Unmarshal(b, &adrs)
	return adrs, err
}
