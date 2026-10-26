package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

// TransactionHandler menangani rute HTTP untuk barang masuk dan keluar
type TransactionHandler struct {
	service *services.TransactionService
}

// NewTransactionHandler membuat instance baru TransactionHandler
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Record menangani pencatatan barang masuk atau keluar (POST /api/spareparts/:id/transactions)
func (h *TransactionHandler) Record(c *gin.Context) {
	idParam := c.Param("id")
	sparepartID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sparepart tidak valid"})
		return
	}

	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Otentikasi dibutuhkan"})
		return
	}
	userID := val.(uint)

	var req services.RecordTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid"})
		return
	}

	transaction, updatedSparepart, err := h.service.RecordTransaction(userID, uint(sparepartID), req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaksi stok berhasil dicatat",
		"transaction": transaction,
		"sparepart":   updatedSparepart,
	})
}

// GetHistory menangani riwayat transaksi sparepart (GET /api/spareparts/:id/transactions)
func (h *TransactionHandler) GetHistory(c *gin.Context) {
	idParam := c.Param("id")
	sparepartID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sparepart tidak valid"})
		return
	}

	history, err := h.service.GetTransactionsBySparepartID(uint(sparepartID))
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}
