package models

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	missingKey      = "Missing key for value '%s'\n"
	missingOperator = "The key '%s' is missing an operator as the last character choose one of %v\n"
	missingValue    = "The key '%s' has no value to check against\n"
	missingValues   = "The key '%s' has no values to check against\n"
)

var validOperators = [5]string{">", ":", "<", ">=", "<="}

type SearchFilter struct {
	Type      string
	DateRange []DateFilter
	Value     map[string]string
	Values    map[string][]string
}

func (sf *SearchFilter) Validate() error {
	sf.defaultType()

	for _, d := range sf.DateRange {
		if err := d.Validate(); err != nil {
			return err
		}
	}

	if err := sf.keyValuesValid(); err != nil {
		return err
	}

	return nil
}

func (sf *SearchFilter) defaultType() {
	sf.Type = "ticket"
}

func (sf *SearchFilter) keyValuesValid() error {
	var bf bytes.Buffer

	for k, v := range sf.Value {
		if len(k) == 0 {
			bf.WriteString(fmt.Sprintf(missingKey, v))
		} else {
			if !keyHasOperator(k) {
				bf.WriteString(fmt.Sprintf(missingOperator, k, validOperators))
			}
		}

		if len(v) == 0 {
			bf.WriteString(fmt.Sprintf(missingValue, k))
		}
	}

	for k, v := range sf.Values {
		if len(k) == 0 {
			bf.WriteString(fmt.Sprintf(missingKey, v))
		} else {
			if !keyHasOperator(k) {
				bf.WriteString(fmt.Sprintf(missingOperator, k, validOperators))
			}
		}

		if len(v) == 0 {
			bf.WriteString(fmt.Sprintf(missingValues, k))
		}
	}

	errMsg := strings.TrimRight(bf.String(), "\n")

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
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
