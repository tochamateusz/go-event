package worker

import (
	"context"
	"os"
	internal_client "tickets/modules/clients"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	external_clients "github.com/ThreeDotsLabs/go-event-driven/common/clients"
)

type Handler struct {
	Handler func(context.Context, Message) error
	Name    string
}

type Worker struct {
	publisher  message.Publisher
	subscriber message.Subscriber
	router     *message.Router
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
			return spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{string(msg.Payload)})
		},
	)

	receiptsClient := internal_client.NewReceiptsClient(clients)
	router.AddNoPublisherHandler("issue-receipt-handler",
		"issue-receipt",
		issuerSub,
		func(msg *message.Message) error {
			return receiptsClient.IssueReceipt(context.Background(), string(msg.Payload))
		},
	)

	return &Worker{
		publisher: publisher,
		router:    router,
	}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		newMsg := message.NewMessage(watermill.NewUUID(), []byte(m.TicketID))
		task, ok := TopicsMap[m.Task]
		if !ok {
			continue
		}
		w.publisher.Publish(task, newMsg)
	}
}

func (w *Worker) Run() error {
	err := w.router.Run(context.Background())
	if err != nil {
		return err
	}
	return nil
}
