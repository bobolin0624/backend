package staging

import "time"

type StagingData struct {
	GroupID   string
	TableName string
	Data      map[string]interface{}

	CreatedAt time.Time
}
