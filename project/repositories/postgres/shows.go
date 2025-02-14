package postgres

import (
	"context"
	"tickets/entities"
	"tickets/repositories"

	"github.com/google/uuid"
	"github.com/gookit/goutil/dump"
	"github.com/jmoiron/sqlx"
)

type PostgresShowsRepository struct {
	db *sqlx.DB
}

func (p *PostgresShowsRepository) getOrCreate(show_id string) entities.Shown {
	shown := entities.Shown{}
	if show_id == "" {
		show_id = uuid.NewString()
	}

	row := p.db.QueryRowx(`
    SELECT show_id, amount FROM shows WHERE show_id=$1;`, show_id,
	)
	if row.Err() != nil {
		return entities.Shown{
			ShowId: show_id,
			Amount: 0,
		}
	}

	err := row.Scan(&shown.ShowId, &shown.Amount)
	if err != nil {
		shown.ShowId = show_id
		return shown
	}

	return shown
}

// Save implements repositories.ShowRepository.
func (p *PostgresShowsRepository) Save(ctx context.Context, show_id string, update func(show entities.Shown) entities.Shown) error {
	show := p.getOrCreate(show_id)
	show = update(show)

	_, err := p.db.Exec(`
INSERT
INTO
	shows
		(
			show_id,
			amount
		)
VALUES
	($1, $2);
`, show.ShowId, show.Amount,
	)

	if err != nil {
		return err
	}

	dump.P(show, err)
	return nil
}

func NewShowsRepository(db *sqlx.DB) repositories.ShowRepository {
	return &PostgresShowsRepository{
		db: db,
	}
}
