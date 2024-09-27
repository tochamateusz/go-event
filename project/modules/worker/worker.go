package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	internal_client "tickets/modules/clients"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	external_clients "github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
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

type EventHeader struct {
	Id          string `json:"id"`
	PublishedAt string `json:"published_at"`
}

type TicketBookingConfirmedEvent struct {
	Header        EventHeader
	TicketId      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type TicketBookingCanceledEvent struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

func LogUUIDMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		defer func() {
			if err != nil {

				logrus.
					WithField("message_uuid", msg.UUID).
					WithField("error", err.Error()).
					Info("Message handling error")
			}
		}()

		logrus.
			WithField("message_uuid", msg.UUID).
			WithField("correlation_id", log.CorrelationIDFromContext(msg.Context())).
			Info("Handling a message")
		return next(msg)
	}
}

func NewWorker(router *message.Router) *Worker {

	logger := watermill.NewStdLogger(false, false)
	logger = logger.With(watermill.LogFields{"context": "worker"})

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	clients, err := external_clients.NewClients(os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			correlationId := log.CorrelationIDFromContext(ctx)
			fmt.Printf("\n\nCLIENT CONTEX: %s\n\n", correlationId)

			req.Header.Set("Correlation-ID", correlationId)
			return nil
		})
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

	refundTickerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "tickets-to-refund",
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(LogUUIDMiddleware)

	spreadsheetsClient := internal_client.NewSpreadsheetsClient(clients)
	router.AddNoPublisherHandler("append-to-tracker-confirmed-handler",
		TopicsMap[TicketBookingConfirmed],
		trackerSub,
		func(msg *message.Message) error {
			payload := TicketBookingCanceledEvent{}
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			ctx := log.ContextWithCorrelationID(msg.Context(), msg.Metadata.Get("correlation_id"))
			msg.SetContext(ctx)

			return spreadsheetsClient.AppendRow(msg.Context(), "tickets-to-print",
				[]string{
					payload.TicketID,
					payload.CustomerEmail,
					payload.Price.Amount,
					payload.Price.Currency})
		},
	)

	router.AddNoPublisherHandler("append-to-tracker-refund-handler",
		TopicsMap[TicketBookingCanceled],
		refundTickerSub,
		func(msg *message.Message) error {
			payload := TicketBookingCanceledEvent{}
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			ctx := log.ContextWithCorrelationID(msg.Context(), msg.Metadata.Get("correlation_id"))
			msg.SetContext(ctx)

			return spreadsheetsClient.AppendRow(msg.Context(), "tickets-to-refund",
				[]string{
					payload.TicketID,
					payload.CustomerEmail,
					payload.Price.Amount,
					payload.Price.Currency})
		},
	)

	receiptsClient := internal_client.NewReceiptsClient(clients)
	router.AddNoPublisherHandler("issue-receipt-handler",
		TopicsMap[TicketBookingConfirmed],
		issuerSub,
		func(msg *message.Message) error {
			payload := TicketBookingConfirmedEvent{}
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			ctx := log.ContextWithCorrelationID(msg.Context(), msg.Metadata.Get("correlation_id"))
			msg.SetContext(ctx)

			return receiptsClient.IssueReceipt(msg.Context(), internal_client.IssueReceiptRequest{
				TicketID: payload.TicketId,
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
