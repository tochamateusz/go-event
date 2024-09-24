package worker

import (
	"context"
	"os"
	internal_client "tickets/modules/clients"

	external_clients "github.com/ThreeDotsLabs/go-event-driven/common/clients"
)

type Worker struct {
	queue              chan Message
	spreadsheetsClient internal_client.SpreadsheetsClient
	receiptsClient     internal_client.ReceiptsClient
}

func NewWorker() *Worker {

	clients, err := external_clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	receiptsClient := internal_client.NewReceiptsClient(clients)
	spreadsheetsClient := internal_client.NewSpreadsheetsClient(clients)

	return &Worker{
		receiptsClient:     receiptsClient,
		spreadsheetsClient: spreadsheetsClient,
		queue:              make(chan Message, 100),
	}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}

func (w *Worker) Run() {
	for msg := range w.queue {
		switch msg.Task {
		case TaskIssueReceipt:
			{
				err := w.receiptsClient.IssueReceipt(context.Background(), msg.TicketID)
				if err != nil {
					w.Send(msg)
				}
			}
		case TaskAppendToTracker:
			err := w.spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{msg.TicketID})
			if err != nil {
				w.Send(msg)
			}
		}
	}
}
