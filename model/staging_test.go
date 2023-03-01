package model

import (
	"testing"
)

func TestStagingCreateValid(t *testing.T) {
	tests := []struct {
		stagingCreate StagingCreate
		expected      bool
	}{
		{
			stagingCreate: StagingCreate{
				Table: "politicians",
				SearchBy: StagingCreateSearchBy{
					"name":     "foo",
					"numumber": 2,
				},
				Fields: StagingCreateFields{
					"name": "bar",
					"party_id": map[string]any{
						"name": "KMT",
					},
					"birthdate": "1990-01-01",
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		if valid := test.stagingCreate.Valid(); valid != test.expected {
			t.Errorf("stagingCreate: %v, Valid(): %v", test.stagingCreate, valid)
		}
	}
}
