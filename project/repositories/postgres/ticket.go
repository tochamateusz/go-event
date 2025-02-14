package postgres

import (
	"context"
	"tickets/entities"
	"tickets/repositories"

	"github.com/jmoiron/sqlx"
)

type PostgresTicketRepository struct {
	db              *sqlx.DB
	existingTickets map[string]struct{}
}

func NewTicketRepository(db *sqlx.DB) repositories.TicketRepository {
	return &PostgresTicketRepository{db, make(map[string]struct{})}
}

// GetAll implements repositories.TicketRepository.
func (p *PostgresTicketRepository) GetAll(context.Context) []entities.Ticket {
	rows, err := p.db.Query(`
SELECT
	ticket_id, price_amount, price_currency, customer_email
FROM
	tickets;
    `)

	tickets := []entities.Ticket{}
	if err != nil {
		return tickets
	}

	for rows.Next() {
		ticket := entities.Ticket{}
		err := rows.Scan(&ticket.TicketID, &ticket.Price.Amount, &ticket.Price.Currency, &ticket.CustomerEmail)
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}

// Delete implements repositories.TicketRepository.
func (p *PostgresTicketRepository) Delete(ctx context.Context, id string) error {
	_, err := p.db.Exec(`DELETE FROM tickets WHERE ticket_id=$1`, id)
	return err
}

// Save implements repositories.TicketRepository.
func (p *PostgresTicketRepository) Save(ctx context.Context, ticket entities.Ticket) error {
	_, ok := p.existingTickets[ticket.TicketID]
	if ok {
		return nil
	}

	_, err := p.db.Exec(`
INSERT
INTO
	tickets
		(
			ticket_id,
			price_amount,
			price_currency,
			customer_email
		)
VALUES
	($1, $2, $3, $4);
`,
		ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)
	if err != nil {
		return err
	}
	p.existingTickets[ticket.TicketID] = struct{}{}

	return nil
}

