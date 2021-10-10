package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type MockAddressManager struct {
	Result map[string][]*Address
}

func (m *MockAddressManager) Find(s string) []*Address {
	return m.Result[s]
}

func TestAddressService_ServeHTTP(t *testing.T) {
	m := &MockAddressManager{Result: map[string][]*Address{
		"one": {
			&Address{
				Line1: "1",
				Line2: "",
				City:  "ONE",
				State: "IL",
				Zip:   "11111",
			},
		},
	}}

	run(t, m, func(addr string) {

		tests := []struct {
			Query string
		}{
			{Query: "one"},
			{Query: "zero"},
		}

		for _, tt := range tests {
			t.Run(tt.Query, func(t *testing.T) {

				adrs, err := requestAPI(addr, tt.Query)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(adrs, m.Result[tt.Query]) {
					j, err := toJson(adrs)
					if err != nil {
						t.Fatal(err)
					}
					t.Fatalf("unexpected result for query [%v]: %v", tt.Query, string(j))
				}
			})
		}
	})
}

func run(t *testing.T, m AddressManager, f func(addr string)) {
	var addr string
	var s *AddressService

	for port := 8000; port <= 9999; port++ {
		addr = fmt.Sprintf("localhost:%d", port)
		s = NewAddressService(addr, m)
		err := s.Start()

		if err == nil {
			t.Logf("server started at: %v", addr)
			break
		}

		if strings.Contains(err.Error(), "address already in use") {
			t.Logf("skip, port in use: %v", err)
			s = nil
			continue
		}

		t.Fatalf("start server error: %v", err)
	}

	if s == nil {
		t.Fatalf("start server error: no port available")
	}

	defer func() {
		if err := s.Stop(); err != nil {
			t.Logf("stop server error: %v", err)
		}

		if err := s.Wait(); err != nil {
			t.Logf("server exit error: %v", err)
		}
	}()
	f(addr)
}

func requestAPI(addr string, query string) ([]*Address, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://%v/address?query=%v", addr, query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status not OK: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("query response: %v", string(b))

	var adrs []*Address
	if err := json.Unmarshal(b, &adrs); err != nil {
		return nil, err
	}
	return adrs, nil
}
