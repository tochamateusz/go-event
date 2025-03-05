package http

import (
	"errors"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type responseMoney struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type ResponseTicket struct {
	TicketID      string        `json:"ticket_id"`
	CustomerEmail string        `json:"customer_email"`
	Price         responseMoney `json:"price"`
}

func From(ticket entities.Ticket) ResponseTicket {
	return ResponseTicket{
		TicketID:      ticket.TicketID,
		CustomerEmail: ticket.CustomerEmail,
		Price: responseMoney{
			Amount:   ticket.Price.Amount,
			Currency: ticket.Price.Currency,
		},
	}
}

func (h Handler) Tickets(c echo.Context) error {
	tickets := h.ticketRepository.GetAll(c.Request().Context())
	ticketsResponse := []ResponseTicket{}
	for _, ticket := range tickets {
		ticketsResponse = append(ticketsResponse, From(ticket))
	}
	return c.JSON(http.StatusOK, ticketsResponse)
}

type ticketsStatusRequest struct {
	Tickets []ticketStatusRequest `json:"tickets"`
}

type ticketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
	BookingID     string         `json:"booking_id"`
}

var IdempotencyKeyMissing = errors.New("Idempotency-Key header missing")

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request ticketsStatusRequest
	idempotencyKeys := c.Request().Header["Idempotency-Key"]
	var idempotencyKey string

	if len(idempotencyKeys) <= 0 {
		idempotencyKey = uuid.NewString()
	}

	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header: entities.NewEventHeaderWithIdempotencyKey(idempotencyKey),

				TicketID:      ticket.TicketID,
				Price:         ticket.Price,
				CustomerEmail: ticket.CustomerEmail,

				BookingID: ticket.BookingID,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingCanceled event: %w", err)
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}

// {
//   "dead_nation_id": "d0b9d5a0-8e1f-4b1a-9f1a-0e8f5e6b9a1a",
//   "number_of_tickets": 100,
//   "start_time": "2021-01-01T00:00:00Z",
//   "title": "The best show ever",
//   "venue": "The best venue ever"
// }

type shownRequest struct {
	DeadNationId    string `json:"dead_nation_id"`
	NumberOfTickets uint16 `json:"number_of_tickets"`
	StartTime       string `json:"start_time"`
	Title           string `json:"title"`
	Venue           string `json:"venue"`
}

func (h Handler) Show(c echo.Context) error {
	var request shownRequest
	err := c.Bind(&request)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	var showE entities.Shown

	h.showRepository.Save(c.Request().Context(), showE.ShowId, func(show entities.Shown) entities.Shown {
		show.Amount += 1
		showE = show
		return show
	})

	return c.JSON(http.StatusCreated, showE)
}

type bookingRequest struct {
	ShowId          string `json:"show_id"`
	NumberOfTickets uint16 `json:"number_of_tickets"`
	CustomerEmail   string `json:"customer_email"`
}

func (h Handler) BookTickets(c echo.Context) error {

	var request bookingRequest
	err := c.Bind(&request)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	booking_id := uuid.NewString()

	var bookingEntity entities.Booking

	err = h.bookingRepository.Save(c.Request().Context(), booking_id, func(booking entities.Booking) entities.Booking {
		bookingEntity = entities.Booking{
			Id:              booking.Id,
			ShowId:          request.ShowId,
			NumberOfTickets: uint(request.NumberOfTickets),
			CustomerEmail:   request.CustomerEmail,
		}
		return bookingEntity
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, booking_id)
	}

	return c.JSON(http.StatusCreated, struct {
		BookingId string `json:"booking_id"`
	}{
		BookingId: bookingEntity.Id,
	})
}
