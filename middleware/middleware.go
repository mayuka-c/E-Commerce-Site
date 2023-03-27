package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/mayuka-c/e-commerce-site/tokens"
)

func Authentication(tokenGenerator *tokens.TokenGenrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Running Authentication Middleware")
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization Header (`token`) not provided"})
			c.Abort()
			return
		}

		claims, err := tokenGenerator.ValidateToken(ClientToken)
		if err != "" {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uuid", claims.UUID)
		c.Next()
	}
}
