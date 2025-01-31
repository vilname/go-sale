package middleware

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"

	"net/http"
)

func AuthenticationMiddleware(c *gin.Context) {
	request := c.Request

	if strings.Contains(request.URL.Path, "/swagger/") {
		c.Next()
		return
	}

	auth := request.Header.Get("Authorization")
	if auth == "" {
		c.Next()
		return
	}

	jwtString := strings.Split(auth, "Bearer ")[1]

	token, _ := jwt.Parse(jwtString, nil)
	if token == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	userUuid := claims["sub"]
	c.Set("userUuid", userUuid)

	c.Next()
}
