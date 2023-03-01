package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/staging"
)

func MountWorkspaceRoutes(rg *gin.RouterGroup) {
	rg.GET("/staging", listStaging)
	rg.POST("/staging/:id", submitStaging)
	rg.POST("/staging", createStaging)
}

func listStaging(c *gin.Context) {
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit > 100 {
		c.Status(http.StatusBadRequest)
	}

	stagingStore := staging.New()
	staging, err := stagingStore.List(c, int(offset), int(limit))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{
		"staging": staging,
	})
}

func submitStaging(c *gin.Context) {
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

func createStaging(c *gin.Context) {
	var body model.StagingCreate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	stagingStore := staging.New()
	if err := stagingStore.Create(c, &body); err == staging.ErrorStagingBadInput {
		c.Status(http.StatusBadRequest)
		return
	} else if err == staging.ErrorStagingNoChange {
		c.Status(http.StatusNotModified)
		return
	} else if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
