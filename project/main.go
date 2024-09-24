package main

import (
	"net/http"
	"tickets/modules/worker"

	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TicketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

type httpError struct {
	Reason string `json:"reason"`
	Method string `json:"method"`
}

func main() {
	log.Init(logrus.InfoLevel)

	w := worker.NewWorker()
	go w.Run()

	e := commonHTTP.NewEcho()

	e.POST("/tickets-confirmation", func(c echo.Context) error {
		var request TicketsConfirmationRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		for _, ticket := range request.Tickets {

			w.Send(worker.Message{
				Task:     worker.TaskIssueReceipt,
				TicketID: ticket,
			})

			w.Send(worker.Message{
				Task:     worker.TaskAppendToTracker,
				TicketID: ticket,
			})
		}

		return c.NoContent(http.StatusOK)
	})

	logrus.Info("Server starting...")

	err := e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
