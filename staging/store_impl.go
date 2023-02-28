package staging

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) Create(ctx context.Context, record *model.StagingDataCreateRecord) error {
	// Check if exist and return id and flag it update. If not flag it create.
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	id := 0
	query, args := createSearchByQuery(record.Table, record.SearchBy)
	if err := conn.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		log.Println(err)
	}

	// Search for fields that needs searching for ids. If not found return failed.
	// Create fields and flags and insert into staging_data if some fields changes
	return errors.New("TODO")
}

func createSearchByQuery(table string, searchBy model.StagingDataSearchBy) (string, []any) {
	where := []string{}
	args := []any{table}
	i := 2
	keys := []string{}
	for k := range searchBy {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		where = append(where, fmt.Sprintf("$%d = $%d", i, i+1))
		args = append(args, k, searchBy[k])
		i += 2
	}
	query := "SELECT id FROM $1 WHERE " + strings.Join(where, " AND ")
	return query, args
}

// TODO refactor
func (s *impl) List(ctx context.Context, offset, limit int) ([]*model.StagingData, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT id, records, created_at, updated_at
		FROM staging_data
		ORDER BY id DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		return nil, err
	}

	stagingData := []*model.StagingData{}
	for rows.Next() {
		var s model.StagingData
		if err := rows.Scan(&s.Id, &s.Records, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		stagingData = append(stagingData, &s)
	}

	return stagingData, nil
}

// TODO refactor
func (s *impl) Submit(ctx context.Context, id int) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if _, err = conn.Exec(ctx, `
		DELETE FROM staging_data
		WHERE id = $1
	`, id); err != nil {
		return err
	}

	// TODO implement the actual submit

	return nil
}
