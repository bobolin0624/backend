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
	where := []string{}
	args := []any{}
	i := 1
	keys := []string{}
	for k := range sc.SearchBy {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		where = append(where, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, sc.SearchBy[k])
		i += 1
	}

	pks := sc.Table.PkNames()
	fieldVars := sc.Table.PkVars()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE "+strings.Join(where, " AND "), strings.Join(pks, ", "), sc.Table)
	return pks, fieldVars.Vars, query, args
}

type StagingCreateSearchBy map[string]any

func (s StagingCreateSearchBy) Valid() bool {
	for _, v := range s {
		switch v.(type) {
		case int:
		case string:
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
			// check if it's nested SearchBy
			return StagingCreateFields(v.(map[string]any)).Valid()
		case int:
			return true
		case string:
			return true
		}
	}

	return false
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
