package api

import (
	"encoding/json"
	"testing"
)

func TestSimpleAddressManager_Find(t *testing.T) {
	j := `
[
  {
    "line1": "Massachusetts Hall",
    "city": "Cambridge",
    "state": "MA",
    "zip": "02138"
  },
  {
    "line1": "3400 N. Charles St.",
    "city": "Baltimore",
    "state": "MD",
    "zip": "21218"
  },
  {
    "line1": "Roosevelt Way NE",
    "city": "Seattle",
    "state": "WA",
    "zip": "98115"
  },
  {
    "line1": "1600 Holloway Ave",
    "city": "San Francisco",
    "state": "CA",
    "zip": "94132"
  },
  {
    "line1": "1600 Holloway Ave",
    "line2": "Suite 10",
    "city": "San Francisco",
    "state": "CA",
    "zip": "94132"
  },
  {
    "line1": "1600 Holloway Ave",
    "line2": "Suite 20",
    "city": "San Francisco",
    "state": "CA",
    "zip": "94132"
  },
  {
    "line1": "500 S State St",
    "city": "Ann Arbor",
    "state": "MI",
    "zip": "48109"
  },
  {
    "line1": "185 Berry St",
    "line2": "Suite 6100",
    "city": "San Francisco",
    "state": "CA",
    "zip": "94107"
  }
]
`

	var adrs []*Address
	if err := json.Unmarshal([]byte(j), &adrs); err != nil {
		t.Fatal(err)
	}

	for _, adr := range adrs {
		b, err := json.Marshal(adr)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("address: %v", string(b))
	}

	sut := NewSimpleAddressManager(adrs)
	for k, v := range sut.index {
		t.Logf("key: %v", k)

		for _, adr := range v {
			t.Logf("address: %+v", *adr)
		}
	}

	tests := []struct {
		Query string
		Count int
	}{
		{Query: "1600", Count: 3},
		{Query: "MD", Count: 1},
		{Query: "6100", Count: 0},
	}

	for _, tt := range tests {
		t.Run(tt.Query, func(t *testing.T) {
			r := sut.Find(tt.Query)
			for _, adr := range r {
				t.Logf("[%v] address: %+v", tt.Query, *adr)
			}

			if tt.Count != len(r) {
				t.Fatalf("expected: %d, actual: %d", tt.Count, len(r))
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	adrs, err := LoadFromFile("testdata/addresses.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(adrs) != 8 {
		t.Fatal("unexpected result")
	}
}
