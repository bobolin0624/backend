package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/candidate"
	"github.com/taiwan-voting-guide/backend/model"
)

func MountCandidate(rg *gin.RouterGroup) {
	rg.POST("/legislator", createCandidateLy)
}

func createCandidateLy(c *gin.Context) {
	var ly model.CandidateLyRepr
	if err := c.BindJSON(&ly); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := candidate.New().Create(c, ly.Model()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
