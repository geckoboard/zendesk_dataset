package conf

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	missingKey      = "Missing key for value '%s'\n"
	missingOperator = "The key '%s' is missing an operator as the last character choose one of %v\n"
	missingValue    = "The key '%s' has no value to check against\n"
	missingValues   = "The key '%s' has no values to check against\n"
)

var validOperators = [5]string{">", ":", "<", ">=", "<="}

// SearchFilter describes a filter query on Zendesk api.
type SearchFilter struct {
	Type      string              `json:"type"`
	DateRange DateFilters         `json:"date_range"`
	Value     map[string]string   `json:"value"`
	Values    map[string][]string `json:"values"`
}

// Validate checks the search filter and returns an error if not.
func (sf *SearchFilter) Validate() error {
	sf.defaultType()

	for _, d := range sf.DateRange {
		if err := d.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (sf *SearchFilter) defaultType() {
	sf.Type = "ticket"
}

func keyHasOperator(key string) bool {
	lastChar := key[len(key)-1:]

	for _, op := range validOperators {
		if lastChar == op {
			return true
		}
	}

	return false
}

// BuildQuery builds a valid zendesk search api and returns it as a string.
func (sf *SearchFilter) BuildQuery(t *time.Time) string {
	var bf bytes.Buffer

	if sf.Type == "" {
		sf.defaultType()
	}

	bf.WriteString(fmt.Sprintf("type:%s ", sf.Type))
	if len(sf.DateRange) != 0 {
		bf.WriteString(sf.DateRange.BuildQuery(t))
		bf.WriteString(" ")
	}

	sf.buildValue(&bf)
	sf.buildValues(&bf)

	return strings.TrimRight(bf.String(), " ")
}

func (sf *SearchFilter) buildValue(bf *bytes.Buffer) {
	// Maps are randomized but we need them ordered at least for the tests.
	keys := []string{}
	for k := range sf.Value {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		bf.WriteString(k)
		val := sf.Value[k]
		if strings.Contains(val, " ") {
			bf.WriteString(`"`)
			bf.WriteString(val)
			bf.WriteString(`"`)
		} else {
			bf.WriteString(val)
		}
		bf.WriteString(" ")
	}

}

func (sf *SearchFilter) buildValues(bf *bytes.Buffer) {
	keys := []string{}
	idx := 0

	for k := range sf.Values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		idx++
		vals := sf.Values[k]

		for i, y := range vals {
			bf.WriteString(k)
			if strings.Contains(y, " ") {
				bf.WriteString(`"`)
				bf.WriteString(y)
				bf.WriteString(`"`)
			} else {
				bf.WriteString(y)
			}
			if i != len(vals)-1 {
				bf.WriteString(" ")
			}
		}

		if idx != len(sf.Values) {
			bf.WriteString(" ")
		}
	}
}
