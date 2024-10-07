package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"tickets/api"
	"tickets/message"
	"tickets/service"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	spreadsheetsService := api.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := api.NewReceiptsServiceClient(apiClients)

	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var schema = `
CREATE TABLE IF NOT EXISTS tickets (
	ticket_id
		UUID PRIMARY KEY,
	price_amount
		DECIMAL(10,2) NOT NULL,
	price_currency
		CHAR(3) NOT NULL,
	customer_email
		VARCHAR(255) NOT NULL
    );
`

	db.MustExec(schema)

	err = service.New(
		redisClient,
		spreadsheetsService,
		receiptsService,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
