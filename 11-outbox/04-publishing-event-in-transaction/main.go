package main

import (
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"
)

func PublishInTx(
	message *message.Message,
	tx *sql.Tx,
	logger watermill.LoggerAdapter,
) error {
	// your code goes here
	return nil
}
