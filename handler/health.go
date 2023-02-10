package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	// log cookit
	session, err := c.Cookie("user_session")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(session)

	c.String(http.StatusOK, "healthy")
}
