package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"tickets/modules/worker"

	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type TicketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

type TicketStatusRequest struct {
	Tickets []worker.Ticket `json:"tickets"`
}

type httpError struct {
	Reason string `json:"reason"`
	Method string `json:"method"`
}

func main() {
	log.Init(logrus.InfoLevel)
	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	w := worker.NewWorker(router)
	g.Go(func() error {
		return w.Run()
	})

	e := commonHTTP.NewEcho()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/tickets-status", func(c echo.Context) error {

		var request TicketStatusRequest
		err = c.Bind(&request)
		if err != nil {
			logger.Error("can't bind", err, make(watermill.LogFields))
			return err
		}

		for _, ticket := range request.Tickets {
			fmt.Printf("Ticket: %+v\n", ticket)
			body := worker.AppendToTrackerPayload{
				TicketId: ticket.ID,
				Price: worker.Money{
					Amount:   ticket.Price.Amount,
					Currency: ticket.Price.Currency,
				},
			}

			data, _ := json.Marshal(body)
			newMessage := message.NewMessage(watermill.NewUUID(), data)
			w.Send(worker.TaskIssueReceipt, newMessage)

			appendToTrackerPayload := worker.AppendToTrackerPayload{
				TicketId:      ticket.ID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			data, _ = json.Marshal(appendToTrackerPayload)
			newMessage = message.NewMessage(watermill.NewUUID(), data)
			w.Send(worker.TaskAppendToTracker, newMessage)
		}

		return c.NoContent(http.StatusOK)
	})

	e.POST("/tickets-confirmation", func(c echo.Context) error {
		var request TicketsConfirmationRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		for _, ticket := range request.Tickets {
			newMessage := message.NewMessage(watermill.NewUUID(), message.Payload(ticket))
			w.Send(worker.TaskIssueReceipt, newMessage)
			w.Send(worker.TaskAppendToTracker, newMessage)
		}

		return c.NoContent(http.StatusOK)
	})

	g.Go(func() error {
		<-router.Running()

		err := e.Start(":8080")
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	g.Go(func() error {
		// Shut down the HTTP server
		<-ctx.Done()
		return e.Shutdown(ctx)
	})

	// Will block until all goroutines finish
	err = g.Wait()
	if err != nil {
		panic(err)
	}

	logrus.Info("Server starting...")

	err = e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
