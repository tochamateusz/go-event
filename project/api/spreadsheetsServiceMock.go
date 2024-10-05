package api

import (
	"context"
	"sync"
	"tickets/message/event"
)

type SpreadsheetsAPIMock struct {
	lock sync.Mutex
	Rows map[string][][]string
}

// AppendRow implements event.SpreadsheetsAPI.
func (s *SpreadsheetsAPIMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Rows == nil {
		s.Rows = make(map[string][][]string)
	}

	s.Rows[sheetName] = append(s.Rows[sheetName], row)

	return nil
}

func NewSpreadsheetsAPIMock() event.SpreadsheetsAPI {
	return &SpreadsheetsAPIMock{
		lock: sync.Mutex{},
		Rows: map[string][][]string{},
	}
}
