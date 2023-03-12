package model

import (
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

func (sc *StagingCreate) Valid() bool {
	return sc.Table.Valid() && sc.SearchBy.Valid() && sc.Fields.Valid()
}

func (sc *StagingCreate) Query() ([]string, []any, string, []any) {
	return searchQuery(sc.Table, sc.SearchBy)
}

type StagingCreateSearchBy map[string]any

func (s StagingCreateSearchBy) Valid() bool {
	for _, v := range s {
		switch v.(type) {
		case int:
		case string:
		case bool:
		default:
			return false
		}
	}

	return true
}

type StagingCreateFields map[string]any

func (f StagingCreateFields) Valid() bool {
	if len(f) == 0 {
		return false
	}
	for _, v := range f {
		switch v.(type) {
		case map[string]any:
			// check if it's nested search
			nestedSearchBy, ok := mapToNestedSearchBy(v.(map[string]any))
			if !ok {
				return false
			}
			return nestedSearchBy.Valid()
		case int:
		case bool:
		case string:
			return true
		}
	}

	return false
}

type StagingCreateNestedSearch struct {
	Table    StagingTable
	SearchBy StagingCreateSearchBy
}

func (ns *StagingCreateNestedSearch) Query() ([]string, []any, string, []any) {
	return searchQuery(ns.Table, ns.SearchBy)
}

func (ns *StagingCreateNestedSearch) Valid() bool {
	return ns.Table.Valid() && ns.SearchBy.Valid()
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
