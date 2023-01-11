package route

import (
	"github.com/gin-gonic/gin"
	"github.com/taiwan-voting-guide/backend/staging"
)

func MountAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/data", listStagingDataHandler)
	rg.POST("/data", submitStagingDataHandler)
}

func listStagingDataHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"stagingData": []StagingDataRepr{
			{
				Id: 2,
				Records: []staging.Record{
					{
						Table: "politicians",
						Record: map[string]interface{}{
							"name":       "邱議瑩",
							"en_name":    "Chiu Yi-Ying",
							"avatar_url": "https://x.gov.tw/pic.jpg",
						},
					},
					{
						Table: "legislators",
						Record: map[string]interface{}{
							"term":         10,
							"session":      4,
							"committee":    "經濟委員會",
							"onboard_date": "2022/05/08",
						},
					},
				},
				CreatedAt: 1673447667,
				UpdatedAt: 1673447667,
			},
			{
				Id: 1,
				Records: []staging.Record{
					{
						Table: "parties",
						Record: map[string]interface{}{
							"id":                  1,
							"name":                "中國國民黨",
							"status":              1,
							"chairman":            "朱立倫",
							"filing_date":         "1989-02-10T00:00:00Z",
							"phone_number":        "(02)87711234",
							"mailing_address":     "臺北市中山區八德路二段232號",
							"established_date":    "1893-11-24T00:00:00Z",
							"main_office_address": "臺北市中山區八德路二段232號",
						},
					},
				},
				CreatedAt: 1673447657,
				UpdatedAt: 1673447657,
			},
		},
	})
}

func submitStagingDataHandler(c *gin.Context) {
	c.JSON(501, gin.H{})
}

type StagingDataRepr struct {
	Id      int              `json:"id"`
	Records []staging.Record `json:"records"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func StagingDataToRepr(s staging.StagingData) StagingDataRepr {
	return StagingDataRepr{
		Id:      s.Id,
		Records: s.Records,

		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
	}
}
