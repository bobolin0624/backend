package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type StagingCreate struct {
	Table    StagingTable          `json:"table"`
	SearchBy StagingCreateSearchBy `json:"searchBy"`
	Fields   StagingCreateFields   `json:"fields"`
}

func (sc *StagingCreate) Valid() (bool, error) {
	if !sc.Table.Valid() {
		return false, errors.New(fmt.Sprintf("invalid table name: %s", sc.Table))
	}

	for k, v := range sc.SearchBy {
		if !sc.Table.isField(k) {
			return false, errors.New(fmt.Sprintf("invalid searchBy key: %s", k))
		}

		switch v.(type) {
		case float64:
		case bool:
		case string:
		default:
			return false, errors.New(fmt.Sprintf("invalid searchBy value: %v", v))

		}

	}

	if len(sc.Fields) == 0 {
		return false, errors.New("fields is empty")
	}

	for k, v := range sc.Fields {
		if !sc.Table.isField(k) {
			return false, errors.New(fmt.Sprintf("invalid fields key: %s", k))
		}

		switch v.(type) {
		case map[string]any:
			// check if it's nested search
			nestedSearchBy, ok := mapToNestedSearchBy(v.(map[string]any))
			if !ok {
				return false, errors.New(fmt.Sprintf("invalid nested searchBy: %v", v))
			}

			return nestedSearchBy.Valid()
		case float64:
		case bool:
		case string:
		default:
			return false, errors.New(fmt.Sprintf("invalid fields value: %v", v))
		}
	}

	return true, nil
}

func (sc *StagingCreate) Query() ([]string, []any, string, []any) {
	return searchQuery(sc.Table, sc.SearchBy)
}

// TODO revisit this
type StagingCreateSearchBy map[string]any

type StagingCreateFields map[string]any

type StagingCreateNestedSearch struct {
	Table    StagingTable
	SearchBy StagingCreateSearchBy
}

func (ns *StagingCreateNestedSearch) Query() ([]string, []any, string, []any) {
	return searchQuery(ns.Table, ns.SearchBy)
}

func (ns *StagingCreateNestedSearch) Valid() (bool, error) {
	if !ns.Table.Valid() {
		return false, errors.New(fmt.Sprintf("invalid nested search table name: %s", ns.Table))
	}

	for k, v := range ns.SearchBy {
		if !ns.Table.isField(k) {
			return false, errors.New(fmt.Sprintf("invalid nested search searchBy key: %s", k))
		}

		switch v.(type) {
		case float64:
		case bool:
		case string:
		default:
			return false, errors.New(fmt.Sprintf("invalid nested search searchBy value: %v", v))
		}
	}

	return true, nil
}

type Staging struct {
	Id     int            `json:"id"`
	Table  StagingTable   `json:"table"`
	Action StagingAction  `json:"action"`
	Fields map[string]any `json:"fields"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s Staging) KeyString() string {
	strs := []string{}
	for _, pk := range s.Table.PkNames() {
		v, ok := s.Fields[pk]
		if !ok {
			panic("Staging.KeyString: missing pk")
		}
		strs = append(strs, fmt.Sprintf("%v", v))
	}

	return strings.Join(strs, "-")
}

const (
	StagingActionCreate StagingAction = "create"
	StagingActionUpdate StagingAction = "update"
)

type StagingAction string

type StagingFieldCompare struct {
	Changed bool `json:"changed"`
	Old     any  `json:"old"`
	New     any  `json:"new"`
}

func searchQuery(table StagingTable, searchBy StagingCreateSearchBy) ([]string, []any, string, []any) {
	where := []string{}
	args := []any{}
	i := 1
	keys := []string{}
	for k := range searchBy {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		where = append(where, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, searchBy[k])
		i += 1
	}

	pks := table.PkNames()
	fieldVars := table.PkVars()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE "+strings.Join(where, " AND "), strings.Join(pks, ", "), table)
	return pks, fieldVars.Vars, query, args
}

func mapToNestedSearchBy(m map[string]any) (StagingCreateNestedSearch, bool) {
	table, ok := m["table"]
	if !ok {
		return StagingCreateNestedSearch{}, false
	}

	searchByMap, ok := m["searchBy"]
	if !ok {
		return StagingCreateNestedSearch{}, false
	}

	searchBy, ok := searchByMap.(map[string]any)
	if !ok {
		return StagingCreateNestedSearch{}, false
	}

	return StagingCreateNestedSearch{
		Table:    StagingTable(table.(string)),
		SearchBy: searchBy,
	}, true
}
