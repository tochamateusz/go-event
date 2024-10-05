package handlers

import (
	"context"
	"fmt"
	"tickets/entities"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type IssueReceiptHandler struct {
	receiptsService event.ReceiptsService
	payload         struct {
		TicketID string         `json:"ticket_id"`
		Price    entities.Money `json:"price"`
	}
}

// Handle implements cqrs.EventHandler.
func (i *IssueReceiptHandler) Handle(ctx context.Context, event any) error {
	evt := event.(*entities.TicketBookingConfirmed)
	log.FromContext(ctx).Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID: evt.TicketID,
		Price:    evt.Price,
	}

	_, err := i.receiptsService.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}

// HandlerName implements cqrs.EventHandler.
func (i *IssueReceiptHandler) HandlerName() string {
	return "IssueReceiptHandler"
}

// NewEvent implements cqrs.EventHandler.
func (i *IssueReceiptHandler) NewEvent() any {
	return i.payload
}

func NewIssueReceiptHandler(receiptsService event.ReceiptsService) cqrs.EventHandler {
	return &IssueReceiptHandler{
		receiptsService: receiptsService,
		payload: struct {
			TicketID string         `json:"ticket_id"`
			Price    entities.Money `json:"price"`
		}{
			TicketID: "",
			Price: entities.Money{
				Amount:   "",
				Currency: "",
			},
		},
	}
}
