package middlewares

import (
	"ecommerce/tokenjwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autherization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims, err := tokenjwt.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.Role != role && role != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not autherized"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := tokenjwt.ValidateToken(tokenString)
		if err != nil || claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this route"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Email)
		c.Next()
	}
}
