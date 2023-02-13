package politician

import (
	"context"
	"log"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (im *impl) Create(ctx context.Context, p *model.PoliticianRepr) (int64, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer conn.Close(ctx)

	var id int64
	err = conn.QueryRow(ctx, `
		INSERT INTO politicians (name, birthdate, avatar_url)
		VALUES ($1, $2, $3)
		RETURNING id
	`, p.Name, p.Birthdate, p.AvatarUrl).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (im *impl) SearchByNameAndBirthdate(ctx context.Context, name, birthdate string) ([]*model.Politician, error) {
	return nil, nil
}
