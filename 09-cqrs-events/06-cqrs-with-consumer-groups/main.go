package main

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewEventProcessor(
	router *message.Router,
	rdb *redis.Client,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) (*cqrs.EventProcessor, error) {
	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return params.EventName, nil
		},
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb,
				ConsumerGroup: "svc-user." + params.HandlerName,
			}, logger)
		},
		Marshaler: marshaler,
		Logger:    logger,
	})
	return eventProcessor, err
}
