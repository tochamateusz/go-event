package api

import (
	"context"
	"tickets/entities"
)

type TicketRepositoryMock struct {
	tickets map[string]entities.Ticket
}

// Delete implements repositories.TicketRepository.
func (t *TicketRepositoryMock) Delete(ctx context.Context, id string) error {
	delete(t.tickets, id)
	return nil
}

// GetAll implements repositories.TicketRepository.
func (t *TicketRepositoryMock) GetAll(context.Context) []entities.Ticket {
	tickets := []entities.Ticket{}
	for _, ticket := range t.tickets {
		tickets = append(tickets, ticket)
	}
	return tickets
}

// Save implements repositories.TicketRepository.
func (t *TicketRepositoryMock) Save(ctx context.Context, ticket entities.Ticket) error {
	t.tickets[ticket.TicketID] = ticket
	return nil
}

func NewTicketRepositoryMock() *TicketRepositoryMock {
	return &TicketRepositoryMock{
		tickets: make(map[string]entities.Ticket),
	}
}
