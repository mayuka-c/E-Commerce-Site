package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/tokens"
)

func Authentication(tokenGenerator *tokens.TokenGenrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := tokenGenerator.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uuid", claims.UUID)
		c.Next()
	}
}
