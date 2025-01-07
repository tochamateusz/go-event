package event

import (
	"context"
	"tickets/entities"
)

func (h Handler) SaveTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	err := h.ticketRepository.Save(ctx, entities.Ticket{
		TicketID: event.TicketID,
		Price: entities.Money{
			Amount:   event.Price.Amount,
			Currency: event.Price.Currency,
		},
		CustomerEmail: event.CustomerEmail,
	})

	return err
}
