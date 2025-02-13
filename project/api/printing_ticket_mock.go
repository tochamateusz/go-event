package api

import (
	"context"
	"tickets/entities"
	"tickets/message/event"
)

type PrintingTicketMock struct{}

// Print implements event.PrintingTicketService.
func (p *PrintingTicketMock) Print(ctx context.Context, request entities.PrintTicketRequest) (entities.PrintTicketResponse, error) {
	return entities.PrintTicketResponse{}, nil
}

func NewPrintingTicketMock() event.PrintingTicketService {
	return &PrintingTicketMock{}
}
