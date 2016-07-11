package zendesk

import "time"

type ticketPayload struct {
	Tickets  []ticket `json:"results"`
	Count    int      `json:"count"`
	NextPage string   `json:"next_page"`
}

type ticket struct {
	ID        int       `json:"id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
