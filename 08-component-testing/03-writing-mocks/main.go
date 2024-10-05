package main

import (
	"context"
	"time"
)

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error)
}

type ReceiptsServiceMock struct {
	IssuedReceipts []IssueReceiptRequest
}

// IssueReceipt implements ReceiptsService.
func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error) {

	r.IssuedReceipts = append(r.IssuedReceipts, IssueReceiptRequest{
		TicketID: request.TicketID,
		Price: Money{
			Amount:   request.Price.Amount,
			Currency: request.Price.Currency,
		},
	})

	return IssueReceiptResponse{
		ReceiptNumber: request.TicketID,
		IssuedAt:      time.Now(),
	}, nil
}

func NewReceiptsServiceMock() *ReceiptsServiceMock {
	return &ReceiptsServiceMock{
		IssuedReceipts: []IssueReceiptRequest{},
	}
}
