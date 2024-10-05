package message

import (
	"tickets/message/decorators"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) (*cqrs.EventBus, error) {
	pub = decorators.CorrelationPublisherDecorator{Publisher: pub}
	return cqrs.NewEventBus(pub, func(eventName string) string {
		return eventName

	}, cqrs.JSONMarshaler{})
}
