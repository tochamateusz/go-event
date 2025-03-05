package postgres

import (
	"context"
	"tickets/entities"
	"tickets/repositories"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresBookingRepository struct {
	db *sqlx.DB
}

func (p *PostgresBookingRepository) getOrCreate(booking_id string) entities.Booking {
	booking := entities.Booking{}
	if booking_id == "" {
		booking_id = uuid.NewString()
	}

	row := p.db.QueryRowx(`
    SELECT id, show_id, number_of_tickets, customer_email FROM bookings WHERE id=$1;`, booking_id)

	if row.Err() != nil {
		return entities.Booking{
			Id:              booking_id,
			ShowId:          "",
			NumberOfTickets: 0,
			CustomerEmail:   "",
		}
	}

	err := row.Scan(&booking.Id, &booking.ShowId, &booking.NumberOfTickets, &booking.CustomerEmail)
	if err != nil {
		return entities.Booking{
			Id:              booking_id,
			ShowId:          "",
			NumberOfTickets: 0,
			CustomerEmail:   "",
		}
	}

	return booking
}

// Save implements repositories.BookingRepository.
func (p *PostgresBookingRepository) Save(ctx context.Context, booking_id string, update func(show entities.Booking) entities.Booking) error {

	booking := p.getOrCreate(booking_id)
	booking = update(booking)

	_, err := p.db.Exec(`
INSERT
INTO
	bookings
		(
      id,
			show_id,
      number_of_tickets,	
      customer_email
		)
VALUES
	($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET 
    show_id=excluded.show_id,
    number_of_tickets=excluded.number_of_tickets,
    customer_email=excluded.customer_email;
`, booking.Id, booking.ShowId, booking.NumberOfTickets, booking.CustomerEmail,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewBookingRepository(db *sqlx.DB) repositories.BookingRepository {
	return &PostgresBookingRepository{
		db: db,
	}
}
