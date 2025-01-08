package main

import "context"

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

type PaymentsHandler struct {
	repo *PaymentsRepository
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{repo: repo}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	return p.repo.SavePaymentTaken(ctx, event)
}

type PaymentsRepository struct {
	payments     []PaymentTaken
	paymentsHash map[string]struct{}
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	return &PaymentsRepository{
		payments:     []PaymentTaken{},
		paymentsHash: make(map[string]struct{}),
	}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	_, ok := p.paymentsHash[event.PaymentID]
	if !ok {
		p.payments = append(p.payments, *event)
		p.paymentsHash[event.PaymentID] = struct{}{}
		return nil
	}
	return nil
}
