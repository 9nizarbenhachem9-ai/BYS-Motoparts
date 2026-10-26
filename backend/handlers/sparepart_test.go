package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sparepart-backend/handlers"
	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

func TestSparepartHandler_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	sparepartService := services.NewSparepartService(nil)
	sparepartHandler := handlers.NewSparepartHandler(sparepartService)

	router.POST("/api/spareparts", sparepartHandler.Create)
	router.GET("/api/spareparts/:id", sparepartHandler.GetByID)

	// 1. Uji ID tidak valid (bukan angka)
	req, _ := http.NewRequest(http.MethodGet, "/api/spareparts/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Diharapkan status 400 Bad Request untuk ID bukan angka, didapat: %d", w.Code)
	}

	// 2. Uji payload kosong / tidak valid
	body, _ := json.Marshal(map[string]interface{}{
		"kode": "",
		"nama": "",
	})
	req, _ = http.NewRequest(http.MethodPost, "/api/spareparts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Diharapkan status 400 Bad Request untuk kode & nama kosong, didapat: %d", w.Code)
	}
}
