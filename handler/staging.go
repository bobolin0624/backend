package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/staging"
)

func MountWorkspaceRoutes(rg *gin.RouterGroup) {
	rg.GET("/staging/:table", listStaging)
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

	table := model.StagingTable(c.Param("table"))
	if !table.Valid() {
		c.Status(http.StatusBadRequest)
		return
	}

	stagingStore := staging.New()
	stagings, err := stagingStore.List(c, table, int(offset), int(limit))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{
		"stagings": stagings,
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
		log.Printf("bad input: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	stagingStore := staging.New()
	if err := stagingStore.Create(c, &body); errors.Is(err, staging.ErrorStagingBadInput) {
		c.Status(http.StatusBadRequest)
		return
	} else if errors.Is(err, staging.ErrorStagingNoChange) {
		c.Status(http.StatusNotModified)
		return
	} else if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
