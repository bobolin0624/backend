package candidate

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct {
}

func (im *impl) Create(ctx context.Context, candidate *model.Candidate) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err := conn.Query(ctx, `
		INSERT INTO candidates (type, term, politician_id, number, elected, party_id, area, vice_president)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, candidate.Type, candidate.Term, candidate.PoliticianId, candidate.Number, candidate.Elected, candidate.PartyId, candidate.Area, candidate.VicePresident); err != nil {
		return err
	}

	return nil
}
