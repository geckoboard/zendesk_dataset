package conf

import (
	"errors"
	"fmt"
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

	// ErrInvalidAttribute is thrown when the metric attribute is invalid.
	ErrInvalidAttribute = fmt.Errorf("The metric attribute is not valid must be one of %s", validMetricAttributes)
	// ErrInvalidUnit is thrown when the metric unit is invalid.
	ErrInvalidUnit = fmt.Errorf("The metric unit is not valid must be one of %s", validMetricUnits)
	// ErrEmptyMetricOptions is thrown when the metric options are all empty.
	ErrEmptyMetricOptions = errors.New("The metric options must be present to be valid for a metric report")
	// ErrEmptyGrouping is thrown when there are no groupings and it is required by the report.
	ErrEmptyGrouping = errors.New("The metric grouping is required when using detailed_metric report")
	// ErrInvalidGroupUnit is thrown when one of the grouping unit are invalid.
	ErrInvalidGroupUnit = errors.New("The metric group unit is invalid must be one of [minute hour]")
	// ErrFromGreaterThanTo is thrown when the group From is greater than To.
	ErrFromGreaterThanTo = errors.New("The metric group 'from' value must not be greater than the 'to' value")
	// ErrFromEqualToTo is thrown when the group From is equal to To.
	ErrFromEqualToTo = errors.New("The metric group 'from' value must not be equal to the 'to' value")
)

// MetricOption describes the options for metric reports
type MetricOption struct {
	Attribute MetricAttribute `yaml:"attribute"`
	Unit      MetricSubMetric `yaml:"unit"`
	Grouping  []MetricGroup   `yaml:"grouping"`
}

// MetricGroup describes how to group ticket metrics. For instance to group
// ticket metrics by 0-1 hours you would specify TimeGroup{Unit: hour, From: 0, To: 1}.
type MetricGroup struct {
	Unit calendarUnit `yaml:"unit"`
	From int          `yaml:"from"`
	To   int          `yaml:"to"`
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

	if m.IsEmpty() {
		return ErrEmptyMetricOptions
	}

	for _, vma := range validMetricAttributes {
		if vma == m.Attribute {
			match = true
			break
		}
	}

	if !match {
		return ErrInvalidAttribute
	}

	match = false
	for _, vmu := range validMetricUnits {
		if vmu == m.Unit {
			match = true
			break
		}
	}

	if !match {
		return ErrInvalidUnit
	}

	return nil
}

// IsEmpty return true if the metric option is initialized with just the default values
// for that data type or false if one of the attributes are not empty.
func (m MetricOption) IsEmpty() bool {
	return m.Attribute == "" && m.Unit == "" && len(m.Grouping) == 0
}

// GroupingValid validates that at least one MetricGroup is present
// if the len is less than 1 it will return an error and checks each
// Metric Group is valid as well.
func (m MetricOption) GroupingValid() error {
	if len(m.Grouping) < 1 {
		return ErrEmptyGrouping
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
			return ErrInvalidGroupUnit
		}
	}

	if m.From > m.To {
		return ErrFromGreaterThanTo
	}

	if m.From == m.To {
		return ErrFromEqualToTo
	}

	return nil
}

// DisplayName returns a string representation of the group name
// and alters the unit into plural version when the To value is
// greater than 1.
func (m MetricGroup) DisplayName() string {
	unit := string(m.Unit)

	if m.To > 1 {
		unit = fmt.Sprintf("%ss", m.Unit)
	}

	return fmt.Sprintf("%d-%d %s", m.From, m.To, unit)
}
