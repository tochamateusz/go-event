package worker

import (
	"context"
	"encoding/json"
	"os"
	internal_client "tickets/modules/clients"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	external_clients "github.com/ThreeDotsLabs/go-event-driven/common/clients"
)

type Handler struct {
	Handler func(context.Context, *message.Message) error
	Name    string
}

type Worker struct {
	publisher  message.Publisher
	subscriber message.Subscriber
	router     *message.Router
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Ticket struct {
	ID            string `json:"ticket_id"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type AppendToTrackerPayload struct {
	TicketId      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

func NewWorker(router *message.Router) *Worker {

	logger := watermill.NewStdLogger(false, false)
	logger = logger.With(watermill.LogFields{"context": "worker"})

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	clients, err := external_clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	trackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)
	if err != nil {
		panic(err)
	}

	issuerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)
	if err != nil {
		panic(err)
	}

	spreadsheetsClient := internal_client.NewSpreadsheetsClient(clients)
	router.AddNoPublisherHandler("append-to-tracker-handler",
		"append-to-tracker",
		trackerSub,
		func(msg *message.Message) error {
			payload := AppendToTrackerPayload{}
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			return spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print",
				[]string{
					payload.TicketId,
					payload.CustomerEmail,
					payload.Price.Amount,
					payload.Price.Currency})
		},
	)

	receiptsClient := internal_client.NewReceiptsClient(clients)
	router.AddNoPublisherHandler("issue-receipt-handler",
		"issue-receipt",
		issuerSub,
		func(msg *message.Message) error {
			payload := internal_client.IssueReceiptRequest{}
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			return receiptsClient.IssueReceipt(context.Background(), internal_client.IssueReceiptRequest{
				TicketID: payload.TicketID,
				Price: internal_client.Money{
					Amount:   payload.Price.Amount,
					Currency: payload.Price.Currency,
				},
			})
		},
	)

	return &Worker{
		publisher: publisher,
		router:    router,
	}
}

func (w *Worker) Send(task Task, msg ...*message.Message) {
	for _, m := range msg {
		task, ok := TopicsMap[task]
		if !ok {
			continue
		}
		w.publisher.Publish(task, m)
	}
}

func (w *Worker) Run() error {
	err := w.router.Run(context.Background())
	if err != nil {
		return err
	}
	return nil
}
