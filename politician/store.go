package politician

import (
	"context"
	"time"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, p *model.PoliticianRepr) (int, error)
	SearchByNameAndBirthdate(ctx context.Context, name string, birthdate *time.Time) ([]*model.Politician, error)
}
