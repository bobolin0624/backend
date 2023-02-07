package handler

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/taiwan-voting-guide/backend/auth"
	"github.com/taiwan-voting-guide/backend/config"
)

func MountAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/google", googleAuthHandler)
}

func googleAuthHandler(c *gin.Context) {
	cookieToken, err := c.Cookie("g_csrf_token")
	if err != nil {
		c.JSON(400, "No CSRF token in Cookie.")
		return
	}

	formToken := c.PostForm("g_csrf_token")
	if formToken == "" {
		c.JSON(400, "No CSRF token in post body.")
		return
	}
	if cookieToken != formToken {
		c.JSON(400, "Failed to verify double submit cookie.")
		return
	}

	credential := c.PostForm("credential")

	authStore := auth.New()
	_, err = authStore.Auth(c, &auth.Info{
		Type: auth.TypeGoogle,
		Google: &auth.InfoGoogle{
			IdToken: credential,
		},
	})
	if err != nil {
		c.JSON(401, err)
	}

	session := sessions.Default(c)

	fmt.Println(session.Get("user_id"))
	session.Set("user_id", credential)
	if err := session.Save(); err != nil {
		c.JSON(500, "Failed to save session.")
		return
	}

	c.SetCookie("session", session.ID(), 3600, "/", config.GetFrontendHost(), false, true)

	c.Redirect(302, config.GetFrontendEndpoint())
}
