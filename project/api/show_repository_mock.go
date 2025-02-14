package api

import (
	"context"
	"tickets/entities"
	"tickets/repositories"
)

type ShowRepositoryMock struct {
}

// Save implements repositories.ShowRepository.
func (s *ShowRepositoryMock) Save(ctx context.Context, show_id string, update func(show entities.Shown) entities.Shown) error {
	panic("unimplemented")
}

func NewShowRepositoryMock() repositories.ShowRepository {
	return &ShowRepositoryMock{}
}
