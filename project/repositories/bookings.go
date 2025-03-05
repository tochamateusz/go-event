package repositories

import (
	"context"
	"tickets/entities"
)

type BookingRepository interface {
	Save(ctx context.Context, show_id string, update func(show entities.Booking) entities.Booking) error
}
