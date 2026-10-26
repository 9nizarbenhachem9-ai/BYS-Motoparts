package handlers

import (
	"net/http"
	"strings"

	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler menangani permintaan HTTP terkait autentikasi
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register menangani endpoint POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid"})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if strings.Contains(err.Error(), "sudah terdaftar") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Akun admin berhasil didaftarkan",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Login menangani endpoint POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid"})
		return
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Me menangani endpoint GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Otentikasi dibutuhkan"})
		return
	}

	userID := val.(uint)
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"created_at": user.CreatedAt,
		},
	})
}

