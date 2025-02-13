package api

import (
	"context"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

type PrintingTicket struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
	printed map[string]struct{}
}

func NewPrintingTicket(clients *clients.Clients) *PrintingTicket {
	if clients == nil {
		panic("NewSpreadsheetsAPIClient: clients is nil")
	}

	return &PrintingTicket{clients: clients}
}

func (p *PrintingTicket) Print(ctx context.Context, request entities.PrintTicketRequest) (
	entities.PrintTicketResponse,
	error,
) {

	resp, err := p.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(
		ctx,
		request.FileID,
		request.Content,
	)

	if resp.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).Infof("file %s already exists", request.FileID)
		return entities.PrintTicketResponse{}, nil
	}

	if err != nil {
		return entities.PrintTicketResponse{}, err
	}

	return entities.PrintTicketResponse{}, nil
}
