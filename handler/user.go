package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/handler/middleware"
	"github.com/taiwan-voting-guide/backend/user"
)

func MountUserRoutes(rg *gin.RouterGroup) {
	rg.Use(middleware.MustAuth())
	rg.GET("/", getUser)
}

func getUser(c *gin.Context) {
	userStore := user.New()
	u, err := userStore.Get(c, c.GetString(middleware.UserIdKey))
	if errors.Is(err, user.ErrUserNotFound) {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u.Repr())
}
