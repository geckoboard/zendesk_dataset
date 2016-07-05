package zendesk

import (
	"fmt"
	"time"

	gb "github.com/geckoboard/geckoboard-go"
	"github.com/jnormington/geckoboard_zendesk/conf"
)

const TicketCount = "ticket_counts"

var timeNow = time.Now()

func HandleReports(c *conf.Config) error {
	for _, r := range c.Zendesk.Reports {
		switch r.Name {
		case TicketCount:
			if err := ticketCount(&r, c); err != nil {
				return fmt.Errorf("Processing dataset report '%s' failed with: %s", r.DataSet, ticketCount(&r, c))
			}
		default:
			return fmt.Errorf("Report name %s was not found", r.Name)
		}
	}

	return nil
}

func ticketCount(r *conf.Report, c *conf.Config) error {
	type GData struct {
		GroupedBy   string `json:"grouped_by"`
		TicketCount int    `json:"ticket_count"`
	}

	client := NewClient(&c.Zendesk.Auth, false)

	var gbData []GData

	if r.GroupBy.Key != "" {
		values := r.Filter.Values[r.GroupBy.Key]
		if len(values) == 0 {
			return fmt.Errorf("Group by values key '%s' returned no values to group by", r.GroupBy.Key)
		}

		delete(r.Filter.Values, r.GroupBy.Key)

		for _, v := range values {
			r.Filter.Values[r.GroupBy.Key] = []string{v}
			tp, err := client.SearchTickets(r.Filter.BuildQuery(&timeNow))

			if err != nil {
				return err
			}

			gbData = append(gbData, GData{GroupedBy: v, TicketCount: tp.Count})
		}
	} else {
		r.GroupBy.Name = "All"
		tp, err := client.SearchTickets(r.Filter.BuildQuery(&timeNow))

		if err != nil {
			return err
		}

		gbData = append(gbData, GData{GroupedBy: r.GroupBy.Name, TicketCount: tp.Count})
	}

	schema := gb.DataSet{
		ID: r.DataSet,
		Fields: gb.Fields{
			"grouped_by": gb.Field{
				Type: gb.StringFieldType,
				Name: r.GroupBy.DisplayName(),
			},
			"ticket_count": gb.Field{
				Type: gb.NumberFieldType,
				Name: "Ticket Count",
			},
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
