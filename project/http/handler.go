package http

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	evtBus                *cqrs.EventBus
	spreadsheetsAPIClient SpreadsheetsAPI
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}
