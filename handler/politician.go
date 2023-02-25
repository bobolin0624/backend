package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/handler/middleware"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
	"github.com/taiwan-voting-guide/backend/politician/question"
)

func MountPolitician(rg *gin.RouterGroup) {
	rg.POST("/", createPolitician)
	rg.GET("/", searchPoliticianByNameAndBirthdate)
	rg.POST("/:politicianId/ask", middleware.MustAuth(), askQuestion)
	rg.GET("/:politicianId/questions", listQuestions)
}

func createPolitician(c *gin.Context) {
	var p model.PoliticianRepr
	if err := c.BindJSON(&p); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err := politician.New().Create(c, &p)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func searchPoliticianByNameAndBirthdate(c *gin.Context) {
	name := c.Query("name")
	birthdate := c.Query("birthdate")

	politicians, err := politician.New().SearchByNameAndBirthdate(c, name, birthdate)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	reprs := make([]*model.PoliticianRepr, len(politicians))
	for i, p := range politicians {
		reprs[i] = p.Repr()
	}

	c.JSON(http.StatusOK, gin.H{
		"politicians": reprs,
	})
}

type AskBody struct {
	Category string `json:"category" binding:"required"`
	Question string `json:"question" binding:"required"`
}

func askQuestion(c *gin.Context) {
	var body AskBody
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	politicianId, err := strconv.ParseInt(c.Param("politicianId"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userId := c.GetString(middleware.UserIdKey)

	q := &model.PoliticianQuestionCreate{
		UserId:       userId,
		PoliticianId: politicianId,
		Category:     body.Category,
		Question:     body.Question,
	}

	questionStore := question.New()
	err = questionStore.Create(c, q)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func listQuestions(c *gin.Context) {
	politicianId, err := strconv.ParseInt(c.Param("politicianId"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit > 100 {
		c.Status(http.StatusBadRequest)
	}

	questions, err := question.New().List(c, politicianId, int(offset), int(limit))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	reprs := []*model.PoliticianQuestionRepr{}
	for _, q := range questions {
		reprs = append(reprs, q.Repr())
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": reprs,
	})
}
