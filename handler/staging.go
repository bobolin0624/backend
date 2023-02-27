package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/staging"
)

func MountWorkspaceRoutes(rg *gin.RouterGroup) {
	rg.GET("/staging", listStagingData)
	rg.POST("/staging/:id", submitStagingData)
}

func listStagingData(c *gin.Context) {
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit > 100 {
		c.Status(http.StatusBadRequest)
	}

	stagingStore := staging.New()
	stagingData, err := stagingStore.List(c, int(offset), int(limit))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	repr := []StagingDataRepr{}
	for _, s := range stagingData {
		repr = append(repr, StagingDataToRepr(*s))
	}

	c.JSON(200, gin.H{
		"stagingData": repr,
	})
}

func submitStagingData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	stagingStore := staging.New()
	if err := stagingStore.Submit(c, int(id)); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

type StagingDataRepr struct {
	Id      int                   `json:"id"`
	Records []model.StagingRecord `json:"records"`

	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func StagingDataToRepr(s model.StagingData) StagingDataRepr {
	return StagingDataRepr{
		Id:      s.Id,
		Records: s.Records,

		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
	}
}
