package event

import (
	"context"
	"tickets/entities"
)

func (h Handler) RemoveTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	err := h.ticketRepository.Delete(ctx, event.TicketID)
	return err
}
