package staging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) Create(ctx context.Context, record *model.StagingCreate) error {
	if !record.Valid() {
		return ErrorStagingBadInput
	}
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	staging := model.Staging{
		Table:  record.Table,
		Action: model.StagingActionCreate,
		Fields: record.Fields,
	}

	// Check if the record exist.
	pks, selects, query, args := record.Query()
	if err = conn.QueryRow(ctx, query, args...).Scan(selects...); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// Insert primary keys to fields if the record exist.
	if err == nil {
		// Mark the record as update.
		staging.Action = model.StagingActionUpdate
		for i, pk := range pks {
			switch selects[i].(type) {
			case *string:
				staging.Fields[pk] = *selects[i].(*string)
			case *int:
				staging.Fields[pk] = *selects[i].(*int)
			default:
				log.Printf("pk: %s, value: %v", pk, selects[i])
				return ErrorStagingBadInput
			}
		}
	}

	// Insert the rest of the fields. If the field is a search pattern, then we search for the primary keys.
	for k, v := range record.Fields {
		switch v.(type) {
		case map[string]any:
			fieldJSON, err := json.Marshal(v)
			if err != nil {
				log.Println(err)
				return ErrorStagingBadInput
			}
			var r model.StagingCreate
			if err := json.Unmarshal(fieldJSON, &r); err != nil {
				log.Println(err)
				return ErrorStagingBadInput
			}

			if !r.Valid() {
				log.Printf("%v: %v\n", k, v)
				return ErrorStagingBadInput
			}

			pks, selects, query, args := r.Query()
			if err := conn.QueryRow(ctx, query, args...).Scan(selects...); errors.Is(err, pgx.ErrNoRows) {
				return ErrorStagingFieldDepNotExist
			} else if err != nil {
				log.Println(err)
				return err
			}

			for _, pk := range pks {
				staging.Fields[k] = pk
			}
		case string:
		case int:
		default:
			log.Printf("%v: %v\n", k, v)
			return ErrorStagingBadInput
		}
	}

	fieldsJSON, err := json.Marshal(staging.Fields)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := conn.Exec(ctx, `
		INSERT INTO staging_data (table_name, action, fields)
		VALUES ($1, $2, $3)
	`, staging.Table, staging.Action, fieldsJSON); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *impl) List(ctx context.Context, table model.StagingTable, offset, limit int) ([]*model.Staging, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT id, table_name, fields, action, created_at, updated_at
		FROM staging_data
		ORDER BY id DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	staging := []*model.Staging{}
	for rows.Next() {
		var s model.Staging
		if err := rows.Scan(&s.Id, &s.Table, &s.Fields, &s.Action, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		staging = append(staging, &s)
	}

	if len(staging) == 0 {
		return []*model.Staging{}, nil
	}

	pks := table.Pks()
	_ = table.Fields()
	/// fields := table.Fields()

	// Generate query for existing records for compare

	conds := []string{}
	args := []any{}
	argsIdx := 1
	for _, s := range staging {
		if s.Action != model.StagingActionUpdate {
			continue
		}

		ands := []string{}
		for _, pk := range pks {
			if _, ok := s.Fields[pk]; !ok {
				log.Printf("pk not found: %s, fields: %v\n", pk, s.Fields)
				return nil, ErrorStagingBadInput
			}

			ands = append(ands, fmt.Sprintf("%s = $%d", pk, argsIdx))
			args = append(args, s.Fields[pk])
			argsIdx++
		}

		conds = append(conds, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
	}

	if len(conds) == 0 {
		return staging, nil
	}

	// query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(fields, ", "), table)
	// query += strings.Join(conds, " OR ")
	// log.Println(query)

	// // combine olds and news records and return result
	// olds := map[string]map[string]any{}
	// rows, err = conn.Query(ctx, query, args...)
	// for rows.Next() {
	// 	args := table.FieldVars()
	// 	if err := rows.Scan(args...); err != nil {
	// 		return nil, err
	// 	}

	// 	pkKey := ""
	// 	for _, pk := range pks {
	// 		pkKey +=

	// }

	return staging, nil
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
