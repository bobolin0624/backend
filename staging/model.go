package staging

import "time"

type StagingData struct {
	GroupID   string                 `json:"groupId"`
	TableName string                 `json:"tableName"`
	Records   map[string]interface{} `json:"records"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
