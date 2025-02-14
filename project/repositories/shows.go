package repositories

import (
	"context"
	"tickets/entities"
)

type ShowRepository interface {
	Save(ctx context.Context, show_id string, update func(show entities.Shown) entities.Shown) error
}
