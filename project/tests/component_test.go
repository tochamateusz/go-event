package tests_test

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"tickets/api"
	"tickets/message"
	"tickets/service"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponent(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	sreadsheetsAPIMock := api.NewSpreadsheetsAPIMock()
	receiptsServiceMock := api.NewReceiptsServiceMock()

	go func() {
		err := service.New(
			redisClient,
			sreadsheetsAPIMock,
			receiptsServiceMock,
		).Run(ctx)

		assert.NoError(t, err)
	}()

	waitForHttpServer(t)
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}
