package zendesk

import (
	"fmt"
	"log"
	"time"

	"github.com/geckoboard/zendesk_dataset/conf"
	gb "github.com/geckoboard/zendesk_dataset/geckoboard"
)

// TicketCount is the supported report template method name.
const (
	TicketCountsReport    = "ticket_counts"
	DetailedMetricsReport = "detailed_metrics"
)

var timeNow = time.Now()

// HandleReports takes a conf.Config and iterates over the Zendesk.Reports
// calling the method based on the Report.Name attribute if any errors
// occurs while processing a report it extracts the error and presents it
// to the user or prints that report was successfull and continues
// with the next report if any.
func HandleReports(c *conf.Config) {
	var err error

	for _, r := range c.Zendesk.Reports {

		switch r.Name {
		case TicketCountsReport:
			err = ticketCount(&r, c)
		case DetailedMetricsReport:
			err = detailedMetrics(&r, c)
		default:
			err = fmt.Errorf("Report name %s was not found", r.Name)
		}

		if err != nil {
			log.Printf("ERRO: Processing report '%s' failed with: %s", r.DataSet, err.Error())
			err = nil
		}

		log.Printf("INFO: Processing report '%s' completed successfully", r.DataSet)
	}
}

func ticketCount(r *conf.Report, c *conf.Config) error {
	type GData struct {
		GroupedBy   string `json:"grouped_by"`
		TicketCount int    `json:"ticket_count"`
	}

	client := newClient(&c.Zendesk.Auth, false)

	var gbData []GData

	if r.GroupBy.Key != "" {
		values := r.Filter.Values[r.GroupBy.Key]
		if len(values) == 0 {
			return fmt.Errorf("Group by values key '%s' returned no values to group by", r.GroupBy.Key)
		}

		delete(r.Filter.Values, r.GroupBy.Key)

		for _, v := range values {
			r.Filter.Values[r.GroupBy.Key] = []string{v}

			tp, err := client.SearchTickets(&Query{Params: r.Filter.BuildQuery(&timeNow)})
			if err != nil {
				return err
			}

			gbData = append(gbData, GData{GroupedBy: v, TicketCount: tp.Count})
		}
	} else {
		r.GroupBy.Name = "All"

		tp, err := client.SearchTickets(&Query{Params: r.Filter.BuildQuery(&timeNow)})
		if err != nil {
			return err
		}

		gbData = append(gbData, GData{GroupedBy: r.GroupBy.Name, TicketCount: tp.Count})
	}

	schema := gb.DataSet{
		ID: r.DataSet,
		Fields: gb.Fields{
			"grouped_by":   gb.Field{Type: gb.StringFieldType, Name: r.GroupBy.DisplayName()},
			"ticket_count": gb.Field{Type: gb.NumberFieldType, Name: "Ticket Count"},
		},
	}

	return pushToGeckoboard(&c.Geckoboard, &schema, gbData)
}

func detailedMetrics(r *conf.Report, c *conf.Config) error {
	if err := r.MetricOptions.Valid(); err != nil {
		return err
	}

	if err := r.MetricOptions.GroupingValid(); err != nil {
		return err
	}

	type MetricData struct {
		Grouping string `json:"grouping"`
		Count    int    `json:"count"`
	}

	client := newClient(&c.Zendesk.Auth, true)
	gbData := make([]MetricData, len(r.MetricOptions.Grouping))

	tm, err := client.TicketMetrics(&Query{Params: r.Filter.BuildQuery(&timeNow)})
	if err != nil {
		return err
	}

	// Group the data as per the user requirements.
	for idx, grp := range r.MetricOptions.Grouping {
		var count int
		d := MetricData{Grouping: grp.DisplayName()}

		for _, t := range tm.Tickets {
			var tMetric int

			switch r.MetricOptions.Unit {
			case conf.BusinessMetric:
				tMetric = t.subTimeMetric(r.MetricOptions.Attribute).Business
			case conf.CalendarMetric:
				tMetric = t.subTimeMetric(r.MetricOptions.Attribute).Calendar
			}

			if tMetric >= grp.FromInMinutes() && tMetric < grp.ToInMinutes() {
				count++
			}
		}

		d.Count = count
		gbData[idx] = d
	}

	schema := gb.DataSet{
		ID: r.DataSet,
		Fields: gb.Fields{
			"grouping": gb.Field{Type: gb.StringFieldType, Name: "Grouping"},
			"count":    gb.Field{Type: gb.NumberFieldType, Name: "Count"},
		},
	}

	return pushToGeckoboard(&c.Geckoboard, &schema, gbData)
}

func pushToGeckoboard(c *conf.Geckoboard, schema *gb.DataSet, data interface{}) error {
	//Create the dataset schema
	gConf := gb.New(gb.Config{
		Key: c.APIKey,
		URL: c.URL,
	})

	err := schema.FindOrCreate(gConf)
	if err != nil {
		return err
	}

	err = schema.SendAll(gConf, data)
	if err != nil {
		return err
	}

	return nil
}
