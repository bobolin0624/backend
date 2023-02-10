package handler

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/auth"
	"github.com/taiwan-voting-guide/backend/config"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/user"
)

func MountAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/google", googleAuthHandler)
	rg.GET("/user", getUser)
}

func googleAuthHandler(c *gin.Context) {
	cookieToken, err := c.Cookie("g_csrf_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "No CSRF token in Cookie.")
		return
	}

	formToken := c.PostForm("g_csrf_token")
	if formToken == "" {
		c.JSON(http.StatusBadRequest, "No CSRF token in post body.")
		return
	}
	if cookieToken != formToken {
		c.JSON(http.StatusBadRequest, "Failed to verify double submit cookie.")
		return
	}

	credential := c.PostForm("credential")

	// get google auth result
	authStore := auth.New()
	result, err := authStore.Auth(c, &model.AuthInfo{
		Type: model.AuthTypeGoogle,
		Google: &model.AuthInfoGoogle{
			IdToken: credential,
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// check if user exist
	userStore := user.New()
	u, err := userStore.GetByAuthResult(c, result)
	if err == user.ErrUserNotFound {
		// create user if not exist
		u, err = userStore.CreateByAuthResult(c, result)
		if err != nil {
			log.Println("***")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", u.Id)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to save session.")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, config.GetFrontendEndpoint())
}

func getUser(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized.")
		return
	}

	userStore := user.New()
	u, err := userStore.Get(c, userId.(string))
	if err == user.ErrUserNotFound {
		c.JSON(http.StatusNotFound, "User not found.")
		return
	}

	c.JSON(http.StatusOK, u.Repr())
}
