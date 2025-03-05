package entities

type Booking struct {
	Id              string `json:"id" db:"id"`
	ShowId          string `json:"show_id" db:"show_id"`
	NumberOfTickets uint   `json:"number_of_tickets" db:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email" db:"customer_email"`
}
