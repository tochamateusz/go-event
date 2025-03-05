package http

import (
	"net/http"
	"tickets/repositories"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	spreadsheetsAPIClient SpreadsheetsAPI,
	ticketRepository repositories.TicketRepository,
	bookingRepository repositories.BookingRepository,
	showRepository repositories.ShowRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:              eventBus,
		spreadsheetsAPIClient: spreadsheetsAPIClient,
		ticketRepository:      ticketRepository,
		bookingRepository:     bookingRepository,
		showRepository:        showRepository,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)
	e.POST("/shows", handler.Show)
	e.GET("/tickets", handler.Tickets)

	e.POST("/book-tickets", handler.BookTickets)

	return e
}
