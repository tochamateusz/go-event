package main

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type FollowRequestSent struct {
	From string `json:"from"`
	To   string `json:"to"`
	EventsCounter
}

var CastError = errors.New("can't cast event")

// Handle implements cqrs.EventHandler.
func (f *FollowRequestSent) Handle(ctx context.Context, event any) error {
	_, ok := event.(*FollowRequestSent)
	if !ok {
		return CastError
	}

	err := f.CountEvent()

	return err
}

// HandlerName implements cqrs.EventHandler.
func (f *FollowRequestSent) HandlerName() string {
	return "FollowRequestSent"
}

// NewEvent implements cqrs.EventHandler.
func (f *FollowRequestSent) NewEvent() any {
	return &(*f)
}

type EventsCounter interface {
	CountEvent() error
}

func NewFollowRequestSentHandler(counter EventsCounter) cqrs.EventHandler {
	return &FollowRequestSent{
		EventsCounter: counter,
		From:          "",
		To:            "",
	}
}
