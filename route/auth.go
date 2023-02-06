package route

import (
	"github.com/gin-gonic/gin"
	"github.com/taiwan-voting-guide/backend/auth"
)

type GoogleAuthPayload struct {
	Credential string `json:"credential"`
	CSRFToken  string `json:"g_csrf_token"`
}

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
	res, err := authStore.Auth(c, &auth.Info{
		Type: auth.TypeGoogle,
		Google: &auth.InfoGoogle{
			IdToken: credential,
		},
	})
	if err != nil {
		c.JSON(401, err)
	}

	c.JSON(302, gin.H{
		"res": res.Google.Payload,
	})

}
