package policy

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, q *model.PoliticianPolicy) error
	Update(ctx context.Context, q *model.PoliticianPolicy) error
	List(ctx context.Context, politicianId int) ([]*model.PoliticianPolicy, error)
}
