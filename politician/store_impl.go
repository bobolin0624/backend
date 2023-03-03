package politician

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (im *impl) Create(ctx context.Context, p *model.PoliticianRepr) (int, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer conn.Close(ctx)

	var id int
	err = conn.QueryRow(ctx, `
		INSERT INTO politicians (name, birthdate, avatar_url, sex)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, p.Name, p.Birthdate, p.AvatarUrl, p.Sex).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

type SearchByNameAndBirthdateParams struct {
	Name  string
	Value interface{}
}

func (im *impl) SearchByNameAndBirthdate(ctx context.Context, name string, birthdate *time.Time) ([]*model.Politician, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	params := []SearchByNameAndBirthdateParams{}
	if name != "" {
		params = append(params, SearchByNameAndBirthdateParams{
			Name:  "name",
			Value: name,
		})
	}
	if birthdate != nil {
		params = append(params, SearchByNameAndBirthdateParams{
			Name:  "birthdate",
			Value: *birthdate,
		})
	}

	where := ""
	if len(params) > 0 {
		where = "WHERE "
		for i, p := range params {
			where += p.Name + " = $" + strconv.FormatInt(int64(i+1), 10)
			if i != len(params)-1 {
				where += " AND "
			}
		}
	}

	args := []any{}
	for _, p := range params {
		args = append(args, p.Value)
	}

	rows, err := conn.Query(ctx, "SELECT id, name, birthdate, avatar_url, sex, created_at, updated_at FROM politicians "+where, args...)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var ps []*model.Politician
	for rows.Next() {
		var p model.Politician
		err = rows.Scan(&p.Id, &p.Name, &p.Birthdate, &p.AvatarUrl, &p.Sex, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		ps = append(ps, &p)
	}

	return ps, nil
}
