package http

import (
	"context"
	"tickets/repositories"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eventBus              *cqrs.EventBus
	ticketRepository      repositories.TicketRepository
	showRepository        repositories.ShowRepository
	spreadsheetsAPIClient SpreadsheetsAPI
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}
