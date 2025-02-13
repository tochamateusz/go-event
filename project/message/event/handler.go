package event

import (
	"context"
	"tickets/entities"
	"tickets/repositories"
)

type Handler struct {
	spreadsheetsService SpreadsheetsAPI
	receiptsService     ReceiptsService
	printTicketService  PrintingTicketService
	ticketRepository    repositories.TicketRepository
}

func NewHandler(
	spreadsheetsService SpreadsheetsAPI,
	receiptsService ReceiptsService,
	printTicketService PrintingTicketService,
	ticketRepository repositories.TicketRepository,
) Handler {
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	if ticketRepository == nil {
		panic("missing ticketRepository")
	}

	return Handler{
		spreadsheetsService: spreadsheetsService,
		receiptsService:     receiptsService,
		printTicketService:  printTicketService,
		ticketRepository:    ticketRepository,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type PrintingTicketService interface {
	Print(ctx context.Context, request entities.PrintTicketRequest) (entities.PrintTicketResponse, error)
}
