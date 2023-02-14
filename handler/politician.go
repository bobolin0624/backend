package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
)

func MountPolitician(rg *gin.RouterGroup) {
	rg.POST("/", createPolitician)
	rg.GET("/", searchPoliticianByNameAndBirthdate)
}

func createPolitician(c *gin.Context) {
	var p model.PoliticianRepr
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(400, "bad request")
		return
	}

	_, err := politician.New().Create(c, &p)
	if err != nil {
		c.JSON(500, "internal server error")
		return
	}

	c.JSON(200, nil)
}

func searchPoliticianByNameAndBirthdate(c *gin.Context) {
	name := c.Query("name")
	birthdate := c.Query("birthdate")

	if name == "" {
		c.JSON(400, "bad request")
		return
	}

	politicians, err := politician.New().SearchByNameAndBirthdate(c, name, birthdate)
	if err != nil {
		c.JSON(500, "internal server error")
		return
	}

	reprs := make([]*model.PoliticianRepr, len(politicians))
	for i, p := range politicians {
		reprs[i] = p.Repr()
	}

	c.JSON(200, gin.H{
		"politicians": reprs,
	})
}
