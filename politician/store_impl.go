package politician

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

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
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	var rows pgx.Rows
	if birthdate == "" {
		rows, err = conn.Query(ctx, `
		SELECT id, name, birthdate, avatar_url
		FROM politicians
		WHERE name = $1

	`, name)
	} else {
		rows, err = conn.Query(ctx, `
		SELECT id, name, birthdate, avatar_url
		FROM politicians
		WHERE name = $1 AND birthdate = $2
	
	`, name, birthdate)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var ps []*model.Politician
	for rows.Next() {
		var p model.Politician
		var t time.Time
		err = rows.Scan(&p.Id, &p.Name, &t, &p.AvatarUrl)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		p.Birthdate = t.Format("2006-01-02")
		ps = append(ps, &p)
	}

	return ps, nil
}
