package event

import (
	"context"
	"fmt"
	"tickets/entities"
)

func (h Handler) TicketBookingConfirmed(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	request := entities.PrintTicketRequest{
		FileID:  fmt.Sprintf("%s-ticket.html", event.TicketID),
		Content: fmt.Sprintf("%s, %s", event.TicketID, event.Price.Amount),
	}
	_, err := h.printTicketService.Print(ctx, request)
	if err != nil {
		return err
	}
	err = h.eventBus.Publish(ctx, entities.TicketPrinted{
		Header:   entities.NewEventHeader(),
		TicketID: event.TicketID,
		FileName: request.FileID,
	})

	return err
}
