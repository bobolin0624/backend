package staging

import (
	"testing"

	"github.com/taiwan-voting-guide/backend/model"
)

func TestCreateSearchByQuery(t *testing.T) {
	tests := []struct {
		table         string
		searchBy      model.StagingDataSearchBy
		expectedQuery string
		expectedArgs  []any
	}{
		{
			table: "party",
			searchBy: model.StagingDataSearchBy{
				"name": "KMT",
			},
			expectedQuery: "SELECT id FROM $1 WHERE $2 = $3",
			expectedArgs:  []any{"party", "name", "KMT"},
		},
		{
			table: "politician",
			searchBy: model.StagingDataSearchBy{
				"name": "蔡英文",
				"num":  2,
			},
			expectedQuery: "SELECT id FROM $1 WHERE $2 = $3 AND $4 = $5",
			expectedArgs:  []any{"politician", "name", "蔡英文", "num", 2},
		},
	}

	for _, test := range tests {
		query, args := createSearchByQuery(test.table, test.searchBy)
		if query != test.expectedQuery {
			t.Errorf("\nquery:\t\t%s\nexpected:\t%s", query, test.expectedQuery)
		}
		for i := range args {
			if args[i] != test.expectedArgs[i] {
				t.Errorf("\nargs:\t\t%v\nexpected:\t%v", args, test.expectedArgs)
			}
		}
	}

}
