package politician

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, p *model.Politician) (int64, error)
	SearchByNameAndBirthdate(ctx context.Context, name, birthdate string) ([]*model.Politician, error)
}
