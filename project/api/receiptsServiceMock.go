package api

import (
	"context"
	"sync"
	"tickets/entities"
	"tickets/message/event"
	"time"
)

type ReceiptsServiceMock struct {
	mock           sync.Mutex
	IssuedReceipts []entities.IssueReceiptRequest
}

// IssueReceipt implements event.ReceiptsService.
func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	r.mock.Lock()
	defer r.mock.Unlock()

	r.IssuedReceipts = append(r.IssuedReceipts, request)

	return entities.IssueReceiptResponse{
		ReceiptNumber: "mocked-receipt-number",
		IssuedAt:      time.Now(),
	}, nil
}

func NewReceiptsServiceMock() event.ReceiptsService {
	return &ReceiptsServiceMock{
		mock:           sync.Mutex{},
		IssuedReceipts: []entities.IssueReceiptRequest{},
	}
}
