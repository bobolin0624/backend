package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/middleware"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
	"github.com/taiwan-voting-guide/backend/politician/question"
)

func MountPolitician(rg *gin.RouterGroup) {
	rg.POST("/", createPolitician)
	rg.GET("/", searchPoliticianByNameAndBirthdate)
	rg.POST("/ask/:politicianId", middleware.MustAuth(), askQuestion)
}

func createPolitician(c *gin.Context) {
	var p model.PoliticianRepr
	if err := c.ShouldBindJSON(&p); err != nil {
		c.Status(400)
		return
	}

	_, err := politician.New().Create(c, &p)
	if err != nil {
		c.Status(500)
		return
	}

	c.Status(201)
}

func searchPoliticianByNameAndBirthdate(c *gin.Context) {
	name := c.Query("name")
	birthdate := c.Query("birthdate")

	politicians, err := politician.New().SearchByNameAndBirthdate(c, name, birthdate)
	if err != nil {
		c.Status(500)
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
		c.Status(400)
		return
	}

	politicianId, err := strconv.ParseInt(c.Param("politicianId"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}

	userId := c.GetString("user_id")

	q := &model.PoliticianQuestionCreate{
		UserId:       userId,
		PoliticianId: politicianId,
		Category:     body.Category,
		Question:     body.Question,
	}

	questionStore := question.New()
	err = questionStore.Create(c, q)
	if err != nil {
		c.Status(500)
		return
	}

	c.Status(201)
}
