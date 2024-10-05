package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewEventProcessor(eventHandlers []cqrs.EventHandler, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) (*cqrs.EventProcessor, error) {

	evtProcessor, err := cqrs.NewEventProcessor(
		eventHandlers,
		func(eventName string) string {
			return eventName
		},
		func(handlerName string) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client: rdb,
			}, watermillLogger)
		},
		cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
		watermillLogger,
	)

	return evtProcessor, err
}
