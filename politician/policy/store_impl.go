package policy

import (
	"context"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

type impl struct{}

func New() Store {
	return &impl{}
}

func (im *impl) Create(ctx context.Context, q *model.PoliticianPolicy) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, `
		INSERT INTO politician_policies (politician_id, category, content)
		VALUES ($1, $2, $3)
	`, q.PoliticianId, q.Category, q.Content); err != nil {
		return err
	}

	return nil
}

func (im *impl) Update(ctx context.Context, q *model.PoliticianPolicy) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, `
		UPDATE politician_policies SET content = $3
		WHERE politician_id = $1 AND category = $2
	`, q.PoliticianId, q.Category, q.Content); err != nil {
		return err
	}

	return nil
}

func (im *impl) List(ctx context.Context, politicianId int) ([]*model.PoliticianPolicy, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT pp.politician_id, pp.category, pp.content
		FROM politician_policies pp
		WHERE politician_id = $1
	`, politicianId)

	if err != nil {
		return nil, err
	}

	policies := []*model.PoliticianPolicy{}
	for rows.Next() {
		q := model.PoliticianPolicy{}
		if err := rows.Scan(&q.PoliticianId, &q.Category, &q.Content); err != nil {
			return nil, err
		}
		policies = append(policies, &q)
	}

	return policies, nil
}
