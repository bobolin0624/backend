package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
)

func MountPolitician(rg *gin.RouterGroup) {
	rg.POST("/", createPolitician)
}

func createPolitician(c *gin.Context) {
	var p model.PoliticianRepr
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(400, "bad request")
		return
	}

	politicianStore := politician.New()
	id, err := politicianStore.Create(c, &p)
	if err != nil {
		c.JSON(500, "internal server error")
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}
