package http

import (
	"errors"
	"fmt"
	"net/http"
	"tickets/entities"

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
	idempotencyKey := c.Request().Header["Idempotency-Key"][0]
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, IdempotencyKeyMissing)
		return IdempotencyKeyMissing
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
