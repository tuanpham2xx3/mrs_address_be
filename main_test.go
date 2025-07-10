package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vietnam-admin-api/handlers"
	"vietnam-admin-api/models"
	"vietnam-admin-api/services"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create test data service
	dataService := services.NewDataService("./testdata")

	// Create test handler
	apiHandler := handlers.NewAPIHandler(dataService, "test")

	// Setup router
	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", apiHandler.Health)
		v1.GET("/provinces", apiHandler.GetProvinces)
		v1.GET("/provinces/:code", apiHandler.GetProvince)
		v1.GET("/wards", apiHandler.GetWards)
		v1.GET("/wards/:code", apiHandler.GetWard)
		v1.GET("/search", apiHandler.GlobalSearch)
	}

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", w.Code)
	}

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}
}

func TestProvincesEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/provinces", nil)
	router.ServeHTTP(w, req)

	// Should return 503 if no test data, or 200 if test data exists
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", w.Code)
	}
}

func TestWardsEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/wards", nil)
	router.ServeHTTP(w, req)

	// Should return 503 if no test data, or 200 if test data exists
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", w.Code)
	}
}

func TestSearchEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=ha", nil)
	router.ServeHTTP(w, req)

	// Should return 503 if no test data, or 200/400 if test data exists
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 200, 400 or 503, got %d", w.Code)
	}
}

// Benchmark tests
func BenchmarkHealthEndpoint(b *testing.B) {
	router := setupTestRouter()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/health", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkProvincesEndpoint(b *testing.B) {
	router := setupTestRouter()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/provinces", nil)
		router.ServeHTTP(w, req)
	}
}
