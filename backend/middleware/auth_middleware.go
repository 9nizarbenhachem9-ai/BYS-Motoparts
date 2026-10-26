package middleware

import (
	"net/http"
	"strings"

	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memverifikasi bearer JWT token pada header Authorization
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token otentikasi tidak ditemukan"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format header Authorization harus Bearer <token>"})
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token otentikasi tidak valid atau telah kedaluwarsa"})
			return
		}

		// Simpan user_id dan username ke dalam context request
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
