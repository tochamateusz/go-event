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
	taskMap    map[Task]Handler
	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewWorker() *Worker {

	logger := watermill.NewStdLogger(false, false)

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

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	spreadsheetsClient := internal_client.NewSpreadsheetsClient(clients)
	receiptsClient := internal_client.NewReceiptsClient(clients)

	taskMap := map[Task]Handler{
		TaskAppendToTracker: {
			Handler: func(ctx context.Context, m Message) error {
				return spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{m.TicketID})
			},
			Name: "append-to-tracker",
		},
		TaskIssueReceipt: {
			Handler: func(ctx context.Context, m Message) error {
				return receiptsClient.IssueReceipt(context.Background(), m.TicketID)
			},
			Name: "issue-receipt",
		},
	}

	return &Worker{
		taskMap:    taskMap,
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		newMsg := message.NewMessage(watermill.NewUUID(), []byte(m.TicketID))
		task, ok := w.taskMap[m.Task]
		if !ok {
			continue
		}
		w.publisher.Publish(task.Name, newMsg)
	}
}

func (w *Worker) Run() {

	for task, handler := range w.taskMap {
		go func(task Task, handler Handler) {
			messages, err := w.subscriber.Subscribe(context.Background(), handler.Name)
			if err != nil {
				panic(err)
			}

			for msg := range messages {
				err = handler.Handler(msg.Context(), Message{
					Task:     task,
					TicketID: string(msg.Payload),
				})
				if err != nil {
					msg.Nack()
				} else {
					msg.Ack()
				}
			}

		}(task, handler)
	}
}
