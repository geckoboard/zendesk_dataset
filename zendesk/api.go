package zendesk

import "time"

// MetricSet describes a subset of metrics under a ticket.
type MetricSet struct {
	ReplyTime          SubTimeMetric `json:"reply_time_in_minutes"`
	FullResolutionTime SubTimeMetric `json:"full_resolution_time_in_minutes"`
}

// SubTimeMetric describe metrics with business and calendar values.
type SubTimeMetric struct {
	Business int `json:"business"`
	Calendar int `json:"calendar"`
}

// TicketPayload the payload returned for search api for type:ticket.
type TicketPayload struct {
	Tickets  []Ticket `json:"results"`
	Count    int      `json:"count"`
	NextPage string   `json:"next_page"`
}

// Ticket makes each Ticket under TicketPayload.
type Ticket struct {
	ID        int       `json:"id"`
	Tags      []string  `json:"tags"`
	Metrics   MetricSet `json:"metric_set"`
	CreatedAt time.Time `json:"created_at"`
}

// TicketMetrics is the tickets/show_many.json schema.
type TicketMetrics struct {
	Tickets  []Ticket `json:tickets`
	Count    int      `json:"count"`
	NextPage string   `json:"next_page"`
}
