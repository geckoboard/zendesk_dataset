package conf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

type calendarUnit string
type dateAttribute string
type DateFilters []DateFilter

const (
	minute calendarUnit = "minute"
	hour   calendarUnit = "hour"
	day    calendarUnit = "day"
	month  calendarUnit = "month"
	year   calendarUnit = "year"

	created dateAttribute = "created"
	updated dateAttribute = "updated"
	solved  dateAttribute = "solved"
	dueDate dateAttribute = "due_date"

	apiDateFormat = "2006-01-02"
)

var validAttributes = [4]dateAttribute{created, updated, solved, dueDate}
var validCalendarUnits = [3]calendarUnit{day, month, year}
var validDateOperators = [3]string{">", ":", "<"}

// DateFilter represents a date filter on the zendesk search api
type DateFilter struct {
	Attribute dateAttribute `json:"attribute"`
	Unit      calendarUnit  `json:"unit"`
	Custom    string        `json:"custom"`
	Past      int           `json:"past"`
}

// Validate returns first error it occurs
func (df *DateFilter) Validate() error {
	df.defaultAttribute()

	if !df.attributeValid() {
		return fmt.Errorf("Attribute is required one of %v", validAttributes)
	}

	if df.Custom != "" && df.Unit != "" && df.Past != 0 {
		return errors.New("Can't use both the unit, past and custom either unit & past or custom on its own")
	}

	if df.Custom == "" && df.Unit == "" && df.Past == 0 {
		return errors.New("Date range or custom input is required")
	}

	if df.Unit != "" && df.Past == 0 {
		return errors.New("Past is required to determine how many units to go into the past from today")
	}

	if !df.unitValid() && df.Past != 0 {
		return fmt.Errorf("Unit is required one of %v", validCalendarUnits)
	}

	if df.Custom != "" && !df.customValid() {
		return errors.New("Custom input requires the operator one of [< : >]")
	}

	return nil
}

func (df *DateFilter) customValid() bool {
	for _, op := range validDateOperators {
		if strings.Index(df.Custom, op) == 0 {
			return true
		}
	}

	return false
}

func (df *DateFilter) defaultAttribute() {
	if df.Attribute == "" {
		df.Attribute = created
	}
}

func (df *DateFilter) unitValid() bool {
	for _, u := range validCalendarUnits {
		if df.Unit == u {
			return true
		}
	}

	return false
}

func (df *DateFilter) attributeValid() bool {
	for _, a := range validAttributes {
		if df.Attribute == a {
			return true
		}
	}

	return false
}

func (df *DateFilter) getDateAPIFormat(t *time.Time) string {
	minusPast := df.Past * -1

	switch df.Unit {
	case day:
		return t.AddDate(0, 0, minusPast).Format(apiDateFormat)
	case month:
		return t.AddDate(0, minusPast, 0).Format(apiDateFormat)
	case year:
		return t.AddDate(minusPast, 0, 0).Format(apiDateFormat)
	}

	return t.Format(apiDateFormat)
}

func (df *DateFilter) BuildQuery(t *time.Time) string {
	var bf bytes.Buffer

	if df.Attribute == "" {
		df.defaultAttribute()
	}

	if t == nil {
		n := time.Now()
		t = &n
	}

	bf.WriteString(string(df.Attribute))
	if df.Custom != "" {
		bf.WriteString(df.Custom)
	} else {
		bf.WriteString(">")
		bf.WriteString(df.getDateAPIFormat(t))
	}

	return bf.String()
}

func (df DateFilters) BuildQuery(t *time.Time) string {
	var bf bytes.Buffer

	for i, d := range df {
		bf.WriteString(d.BuildQuery(t))
		if i != len(df)-1 {
			bf.WriteString(" ")
		}
	}

	return bf.String()
}
