package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/gookit/goutil/dump"
)

func (h Handler) TicketBookingConfirmed(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	request := entities.PrintTicketRequest{
		FileID:  fmt.Sprintf("%s-ticket.html", event.TicketID),
		Content: fmt.Sprintf("%s, %s", event.TicketID, event.Price.Amount),
	}
	dump.P(request)
	_, err := h.printTicketService.Print(ctx, request)

	return err
}
