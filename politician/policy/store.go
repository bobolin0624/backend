package policy

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
)

type Store interface {
	Create(ctx context.Context, q *model.PoliticianPolicyCreate) error
	List(ctx context.Context, politicianId int) ([]*model.PoliticianPolicy, error)
	Search(ctx context.Context, politicianId int, typ string) ([]*model.PoliticianPolicy, error)
}
