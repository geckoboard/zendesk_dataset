package models

import (
	"errors"
	"fmt"
	"strings"
)

type calendarUnit string

const (
	minute calendarUnit = "minute"
	hour   calendarUnit = "hour"
	day    calendarUnit = "day"
	month  calendarUnit = "month"
	year   calendarUnit = "year"
)

type dateAttribute string

const (
	created dateAttribute = "created"
	updated dateAttribute = "updated"
	solved  dateAttribute = "solved"
	dueDate dateAttribute = "due_date"
)

var validAttributes = [4]dateAttribute{created, updated, solved, dueDate}
var validCalendarUnits = [5]calendarUnit{minute, hour, day, month, year}
var validDateOperators = [3]string{">", ":", "<"}

// DateFilter represents a date filter on the zendesk search api
type DateFilter struct {
	Attribute dateAttribute
	Unit      calendarUnit
	Custom    string
	Past      int
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
