package handlers

import (
	"context"
	"tickets/entities"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type AppendToTrackerHandler struct {
	spreadsheetsService event.SpreadsheetsAPI
	payload             entities.Ticket
}

// Handle implements cqrs.EventHandler.
func (a *AppendToTrackerHandler) Handle(ctx context.Context, event any) error {
	evt := event.(*entities.Ticket)

	log.FromContext(ctx).Info("Appending ticket to the tracker")

	return a.spreadsheetsService.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{evt.TicketID, evt.CustomerEmail, evt.Price.Amount, evt.Price.Currency},
	)
}

// HandlerName implements cqrs.EventHandler.
func (a *AppendToTrackerHandler) HandlerName() string {
	return "AppendToTrackerHandler"
}

// NewEvent implements cqrs.EventHandler.
func (a *AppendToTrackerHandler) NewEvent() any {
	return a.payload
}

func NewAppendToTrackerHandler(spreadsheetsService event.SpreadsheetsAPI) cqrs.EventHandler {
	return &AppendToTrackerHandler{
		spreadsheetsService: spreadsheetsService,
		payload: entities.Ticket{
			TicketID: "",
			Price: entities.Money{
				Amount:   "",
				Currency: "",
			},
			CustomerEmail: "",
		},
	}
}
