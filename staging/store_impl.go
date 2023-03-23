package staging

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/pg"
)

func New() Store {
	return &impl{}
}

type impl struct{}

func (s *impl) Create(ctx context.Context, staging model.Staging) error {
	if ok, err := staging.Valid(); !ok {
		return err
	}

	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	searchByJSON, err := json.Marshal(staging.SearchBy)
	if err != nil {
		return err
	}

	fields := model.StagingFields{}

	// Insert the rest of the fields. If the field is a search pattern, then we search for the primary keys.
	for k, v := range staging.Fields {
		switch v.(type) {
		case map[string]any:
			log.Println("nested search")
			fieldJSON, err := json.Marshal(v)
			if err != nil {
				return err
			}

			var ns model.StagingNestedSearch
			if err := json.Unmarshal(fieldJSON, &ns); err != nil {
				return err
			}

			if ok, err := ns.Valid(); !ok {
				return err
			}

			pks, selects, query, args := ns.Query()
			if err := conn.QueryRow(ctx, query, args...).Scan(selects...); errors.Is(err, pgx.ErrNoRows) {
				return ErrorStagingFieldDepNotExist
			} else if err != nil {
				return err
			}

			if len(pks) == 1 {
				fields[k] = sqlVarToAny(selects[0])
				continue
			}

			m := model.StagingFields{}
			for i, pk := range pks {
				m[pk] = sqlVarToAny(selects[i])
			}
			fields[k] = m

		case string:
			fields[k] = v
		case float64:
			fields[k] = v
		case bool:
			fields[k] = v
		default:
			return ErrorStagingBadInput
		}
	}

	fmt.Println(fields)

	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	if _, err := conn.Exec(ctx, `
		INSERT INTO staging_data (table_name, search_by, fields)
		VALUES ($1, $2, $3)
		ON CONFLICT (table_name, search_by, fields)
		DO UPDATE SET updated_at = NOW()
	`, staging.Table, searchByJSON, fieldsJSON); err != nil {
		return err
	}

	return nil
}

func sqlVarToAny(v any) any {
	switch v.(type) {
	case *sql.NullString:
		return v.(*sql.NullString).String
	case *sql.NullInt64:
		return v.(*sql.NullInt64).Int64
	case *sql.NullBool:
		return v.(*sql.NullBool).Bool
	case *sql.NullTime:
		return v.(*sql.NullTime).Time
	default:
		return nil
	}
}

func (s *impl) List(ctx context.Context, table model.StagingTable, offset, limit int) ([]model.StagingResult, error) {
	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(ctx)

	// Query from staging_data
	rows, err := conn.Query(ctx, `
		SELECT id, table_name, search_by, fields, created_at 
		FROM staging_data
		WHERE table_name = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`, table, offset, limit)
	if err != nil {
		return nil, err
	}

	stagings := []*model.Staging{}
	for rows.Next() {
		var s model.Staging
		if err := rows.Scan(&s.Id, &s.Table, &s.SearchBy, &s.Fields, &s.CreatedAt); err != nil {
			return nil, err
		}

		stagings = append(stagings, &s)
	}

	if len(stagings) == 0 {
		return []model.StagingResult{}, nil
	}

	// dedup searchBys
	searchBys := []model.StagingFields{}
	for _, staging := range stagings {
		for _, searchBy := range searchBys {
			if searchBy.Equal(staging.SearchBy) {
				break
			}
		}
		searchBys = append(searchBys, staging.SearchBy)
	}

	// create select condition by searchBys
	conds := []string{}
	args := []any{}
	argsIdx := 1
	for _, searchBy := range searchBys {
		ands := []string{}
		for field, value := range searchBy {
			ands = append(ands, fmt.Sprintf("%s = $%d", field, argsIdx))
			args = append(args, value)
			argsIdx++
		}
		conds = append(conds, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
	}

	if len(conds) == 0 {
		return []model.StagingResult{}, nil
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE ", strings.Join(table.FieldNames(), ", "), table)
	query += strings.Join(conds, " OR ")

	oldFieldsList := []model.StagingFields{}
	rows, err = conn.Query(ctx, query, args...)
	for rows.Next() {
		fieldVars := table.FieldVars()
		if err := rows.Scan(fieldVars.Vars...); err != nil {
			log.Println(err)
			return nil, err
		}

		oldFieldsList = append(oldFieldsList, fieldVars.Map())
	}

	// Combine old and new records and return result
	results := []model.StagingResult{}
	for _, s := range stagings {
		foundFields := []model.StagingFields{}
		for _, oldFields := range oldFieldsList {
			if s.SearchBy.ExistIn(oldFields) {
				foundFields = append(foundFields, oldFields)
			}
		}

		switch len(foundFields) {
		case 0:
			resultFields := []model.StagingResultField{}
			for _, fn := range table.FieldNames() {
				v := s.Fields[fn]
				resultFields = append(resultFields, model.StagingResultField{
					Type:  model.StagingResultFieldTypeValue,
					Name:  fn,
					Value: v,
				})
			}

			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: model.StagingResultStatusCreate,
			})
		case 1:
			// Inject primary key from search result
			oldFields := foundFields[0]
			for _, pk := range table.PkNames() {
				s.Fields[pk] = oldFields[pk]
			}

			resultFields := []model.StagingResultField{}
			for _, fieldName := range table.FieldNames() {
				newVal, newOk := s.Fields[fieldName]
				oldVal, oldOk := oldFields[fieldName]

				compare := model.StagingFieldCompare{}
				if newOk {
					compare.New = newVal
				}
				if oldOk {
					compare.Old = oldVal
				}

				if fieldChanged(oldVal, newVal) {
					compare.Changed = true
				}

				resultFields = append(resultFields, model.StagingResultField{
					Type:  model.StagingResultFieldTypeCompare,
					Name:  fieldName,
					Value: compare,
				})
			}

			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: model.StagingResultStatusUpdate,
			})
		default:
			resultFields := []model.StagingResultField{}
			for _, fn := range table.FieldNames() {
				if v, ok := s.Fields[fn]; ok {
					resultFields = append(resultFields, model.StagingResultField{
						Type:  model.StagingResultFieldTypeValue,
						Name:  fn,
						Value: v,
					})
				}
			}
			results = append(results, model.StagingResult{
				Id:     s.Id,
				Fields: resultFields,
				Status: model.StagingResultStatusConflict,
			})
		}

	}

	return results, nil
}

func fieldChanged(old, new any) bool {
	if new == nil {
		return false
	}

	switch o := old.(type) {
	case int64:
		switch n := new.(type) {
		case int64:
			return o != n
		case float64:
			return float64(o) != n
		}
	case float64:
		switch n := new.(type) {
		case int64:
			return o != float64(n)
		case float64:
			return o != n
		}
	case string:
		n, ok := new.(string)
		if !ok {
			return true
		}
		return o != n
	case bool:
		n, ok := new.(bool)
		if !ok {
			return true
		}
		return o != n
	}

	return true
}

func (s *impl) Submit(ctx context.Context, submit model.StagingSubmit) error {
	conn, err := pg.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	tx.Exec(ctx, fmt.Sprintf("LOCK TABLE %s IN EXCLUSIVE MODE", submit.Table))

	if _, err = tx.Exec(ctx, `
		DELETE FROM staging_data
		WHERE id = $1
	`, submit.Id); err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}
