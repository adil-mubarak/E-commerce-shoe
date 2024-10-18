package middlewares

import (
	"ecommerce/database"
	"ecommerce/models"
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

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := tokenjwt.ValidateToken(tokenString)
		if err != nil {
			if err.Error() == "Token is expired" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired, please refresh"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

			}
			c.Abort()
			return
		}

		if claims.Role != role && role != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
			c.Abort()
			return
		}

		userID := claims.UserID

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.Banned {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is banned"})
			c.Abort()
			return
		}

		c.Set("user_id", claims)
		c.Next()

	}
}
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := tokenjwt.ValidateToken(tokenString)
		if err != nil {
			if err.Error() == "Token is expired" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired, please refresh"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden"})
			c.Abort()
			return
		}

		c.Next()

	}
}
