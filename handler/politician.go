package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
	"github.com/taiwan-voting-guide/backend/politician/question"
)

func MountPolitician(rg *gin.RouterGroup) {
	rg.POST("/", createPolitician)
	rg.GET("/", searchPoliticianByNameAndBirthdate)
	rg.POST("/ask/:politicianId", askQuestion)
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

type AskBody struct {
	Category string `json:"category"`
	Question string `json:"question"`
}

func askQuestion(c *gin.Context) {
	var body AskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, "bad request")
		return
	}

	politicianId, err := strconv.ParseInt(c.Param("politicianId"), 10, 64)
	if err != nil {
		c.JSON(400, "bad request")
		return
	}

	q := &model.PoliticianQuestionCreate{
		UserId:       "TODO",
		PoliticianId: politicianId,
		Category:     body.Category,
		Question:     body.Question,
	}

	questionStore := question.New()
	err = questionStore.Create(c, q)
	if err != nil {
		c.JSON(500, "internal server error")
		return
	}

	c.JSON(200, nil)
}
