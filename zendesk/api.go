package zendesk

import "time"

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
	CreatedAt time.Time `json:"created_at"`
}
