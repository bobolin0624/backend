package politician

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type impl struct{}

func New() Store {
	return &impl{}
}

func (im *impl) Create(ctx context.Context, p *model.Politician) (int64, error) {
	return 0, nil
}

func (im *impl) SearchByNameAndBirthdate(ctx context.Context, name, birthdate string) ([]*model.Politician, error) {
	return nil, nil
}
