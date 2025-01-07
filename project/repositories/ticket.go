package repositories

import (
	"context"
	"tickets/entities"
)

type TicketRepository interface {
	Save(context.Context, entities.Ticket) error
	GetAll(context.Context) []entities.Ticket
	Delete(ctx context.Context, id string) error
}
