package conf

import (
	"errors"
	"fmt"
	"reflect"
)

// MetricAttribute defines the allowed metric attributes.
type MetricAttribute string

// MetricSubMetric is the unit to use from a metric attribute.
type MetricSubMetric string

const (
	//ReplyTime is the first reply time in minutes.
	ReplyTime MetricAttribute = "reply_time"
	//FirstResolutionTime is the first resolution time in minutes.
	FirstResolutionTime MetricAttribute = "first_resolution_time"
	//FullResolutionTime is the full resolution time in minutes.
	FullResolutionTime MetricAttribute = "full_resolution_time"
	//AgentWaitTime is the first agent wait time in minutes.
	AgentWaitTime MetricAttribute = "agent_wait_time"
	//RequesterWaitTime is the first requester wait time in minutes.
	RequesterWaitTime MetricAttribute = "requester_wait_time"
	//OnHoldTime is the first on hold time in minutes.
	OnHoldTime MetricAttribute = "on_hold_time"

	//BusinessMetric relates to the business time.
	BusinessMetric MetricSubMetric = "business"
	//CalendarMetric relates to the calendar time.
	CalendarMetric MetricSubMetric = "calendar"
)

var (
	validMetricAttributes = [6]MetricAttribute{
		ReplyTime,
		FirstResolutionTime,
		FullResolutionTime,
		AgentWaitTime,
		RequesterWaitTime,
		OnHoldTime,
	}

	validMetricUnits = [2]MetricSubMetric{BusinessMetric, CalendarMetric}

	errEmptyMetricOptions = fmt.Errorf("The metric options must be present to be valid for a metric report")
	errInvalidAttribute   = fmt.Errorf("The metric attribute is not valid must be one of %s", validMetricAttributes)
	errInvalidUnit        = fmt.Errorf("The metric unit is not valid must be one of %s", validMetricUnits)
	errEmptyGrouping      = errors.New("The metric grouping is required when using detailed_metric report")
	errInvalidGroupUnit   = errors.New("The metric group unit is invalid must be one of [minute hour]")
	errFromGreaterThanTo  = errors.New("The metric group 'from' value must not be greater than the 'to' value")
	errFromEqualToTo      = errors.New("The metric group 'from' value must not be equal to the 'to' value")
)

// MetricOption describes the options for metric reports
type MetricOption struct {
	Attribute MetricAttribute
	Unit      MetricSubMetric
	Grouping  []MetricGroup
}

// MetricGroup describes how to group ticket metrics. For instance to group
// ticket metrics by 0-1 hours you would specify TimeGroup{Unit: hour, From: 0, To: 1}.
type MetricGroup struct {
	Unit calendarUnit `json:"unit"`
	From int          `json:"from"`
	To   int          `json:"to"`
}

// FromInMinutes returns from converted into minutes based on the unit
// else returns -1 if the unit is not valid.
func (m MetricGroup) FromInMinutes() int {
	return getUnitInMinutes(m.Unit, m.From)
}

// ToInMinutes returns from converted into minutes based on the unit
// else returns -1 if the unit is not valid.
func (m MetricGroup) ToInMinutes() int {
	return getUnitInMinutes(m.Unit, m.To)
}

func getUnitInMinutes(unit calendarUnit, val int) int {
	switch unit {
	case minute:
		return val
	case hour:
		return val * 60
	}

	return -1
}

// Valid return either true or false validating options required for metric reports.
// It doesn't validate against the grouping as not all metric reports require it.
func (m MetricOption) Valid() error {
	var match bool

	if reflect.DeepEqual(m, MetricOption{}) {
		return errEmptyMetricOptions
	}

	for _, vma := range validMetricAttributes {
		if vma == m.Attribute {
			match = true
		}
	}

	if !match {
		return errInvalidAttribute
	}

	match = false
	for _, vmu := range validMetricUnits {
		if vmu == m.Unit {
			match = true
		}
	}

	if !match {
		return errInvalidUnit
	}

	return nil
}

// GroupingValid validates that at least one MetricGroup is present
// if the len is less than 1 it will return an error and checks each
// Metric Group is valid as well.
func (m MetricOption) GroupingValid() error {
	if len(m.Grouping) < 1 {
		return errEmptyGrouping
	}

	for _, g := range m.Grouping {
		if err := g.valid(); err != nil {
			return err
		}
	}

	return nil
}

func (m MetricGroup) valid() error {
	if m.Unit != minute {
		if m.Unit != hour {
			return errInvalidGroupUnit
		}
	}

	if m.From > m.To {
		return errFromGreaterThanTo
	}

	if m.From == m.To {
		return errFromEqualToTo
	}

	return nil
}
