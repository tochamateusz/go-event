package service

import (
	"context"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"
	"tickets/message/event/handlers"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Init(logrus.InfoLevel)
}

type Service struct {
	evtProcessor *cqrs.EventProcessor
	evtBus       *cqrs.EventBus
	echoRouter   *echo.Echo
}

func New(
	redisClient *redis.Client,
	spreadsheetsService event.SpreadsheetsAPI,
	receiptsService event.ReceiptsService,
) Service {
	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	redisPublisher := message.NewRedisPublisher(redisClient, watermillLogger)

	hs := []cqrs.EventHandler{
		handlers.NewIssueReceiptHandler(receiptsService),
		handlers.NewAppendToTrackerHandler(spreadsheetsService),
	}

	evtProcessor, _ := message.NewEventProcessor(hs, redisClient, watermillLogger)
	evtBus, _ := message.NewEventBus(redisPublisher)

	echoRouter := ticketsHttp.NewHttpRouter(
		evtBus,
		spreadsheetsService,
	)

	return Service{
		evtProcessor,
		evtBus,
		echoRouter,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && err != stdHTTP.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
