package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		fmt.Printf("CreateUserAccount: Encounter a problem: %s\n", err.Error())
		return err
	}

	go retry(func() error {
		return h.newsletterClient.AddToNewsletter(u)
	})

	go retry(func() error {
		return h.notificationsClient.SendNotification(u)
	})

	return nil
}

func retry(fnc func() error) {
	for {
		err := fnc()
		if err != nil {
			filename, line := runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).FileLine(runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Entry())
			fmt.Printf("%s, [%d]: Encounter a problem: %s\n", filename, line, err.Error())
			time.Sleep(10 * time.Microsecond)
			continue
		}
		break
	}
}
