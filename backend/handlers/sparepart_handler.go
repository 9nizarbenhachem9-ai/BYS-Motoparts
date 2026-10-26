package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

// SparepartHandler menangani permintaan HTTP terkait pengelolaan sparepart
type SparepartHandler struct {
	service *services.SparepartService
}

// NewSparepartHandler membuat instance baru SparepartHandler
func NewSparepartHandler(service *services.SparepartService) *SparepartHandler {
	return &SparepartHandler{service: service}
}

// Create menangani penambahan sparepart baru (POST /api/spareparts)
func (h *SparepartHandler) Create(c *gin.Context) {
	var req services.CreateSparepartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid"})
		return
	}

	item, err := h.service.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "sudah digunakan") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sparepart berhasil ditambahkan",
		"data":    item,
	})
}

// GetAll menangani pengambilan daftar sparepart dengan pencarian & paginasi (GET /api/spareparts)
func (h *SparepartHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	kategori := c.Query("kategori")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, err := h.service.GetAll(search, kategori, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID menangani pengambilan detail sparepart (GET /api/spareparts/:id)
func (h *SparepartHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sparepart tidak valid"})
		return
	}

	item, err := h.service.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Update menangani pengubahan data sparepart (PUT /api/spareparts/:id)
func (h *SparepartHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sparepart tidak valid"})
		return
	}

	var req services.UpdateSparepartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid"})
		return
	}

	item, err := h.service.Update(uint(id), req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sudah digunakan") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data sparepart berhasil diperbarui",
		"data":    item,
	})
}

// Delete menangani penghapusan sparepart (DELETE /api/spareparts/:id)
func (h *SparepartHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sparepart tidak valid"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sparepart berhasil dihapus",
	})
}
