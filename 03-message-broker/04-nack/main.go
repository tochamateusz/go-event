package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		state := string(msg.Payload)
		fmt.Printf("STATE: %+v\n", state)

		if state == "1" {
			err := alarmClient.StartAlarm()
			if err != nil {
				msg.Nack()
				continue
			}
			msg.Ack()
			continue
		}

		if state == "0" {
			err := alarmClient.StopAlarm()
			if err != nil {
				msg.Nack()
				continue
			}
			msg.Ack()
			continue
		}

	}
}
