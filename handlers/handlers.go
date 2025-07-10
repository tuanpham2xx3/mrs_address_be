package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vietnam-admin-api/models"
	"vietnam-admin-api/services"

	"github.com/gin-gonic/gin"
)

// APIHandler contains the data service and handles HTTP requests
type APIHandler struct {
	dataService *services.DataService
	startTime   time.Time
	version     string
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(dataService *services.DataService, version string) *APIHandler {
	return &APIHandler{
		dataService: dataService,
		startTime:   time.Now(),
		version:     version,
	}
}

// Helper functions

func (h *APIHandler) parseQueryParams(c *gin.Context) (search, typeFilter string, limit, offset int) {
	search = strings.TrimSpace(c.Query("search"))
	typeFilter = strings.TrimSpace(c.Query("type"))

	limit = 50 // default
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	offset = 0 // default
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return
}

func (h *APIHandler) respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, models.APIResponse{
		Success: false,
		Message: message,
	})
}

func (h *APIHandler) checkDataLoaded(c *gin.Context) bool {
	if !h.dataService.IsDataLoaded() {
		h.respondWithError(c, http.StatusServiceUnavailable, "Data not loaded")
		return false
	}
	return true
}

// Province Handlers

// GetProvinces handles GET /api/v1/provinces
func (h *APIHandler) GetProvinces(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	search, typeFilter, limit, offset := h.parseQueryParams(c)

	provinces, total := h.dataService.SearchProvinces(search, typeFilter, limit, offset)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    provinces,
		Pagination: models.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Pages:  int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// GetProvince handles GET /api/v1/provinces/:code
func (h *APIHandler) GetProvince(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	code := c.Param("code")
	if code == "" {
		h.respondWithError(c, http.StatusBadRequest, "Province code is required")
		return
	}

	province, err := h.dataService.GetProvince(code)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, "Province not found")
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    province,
	})
}

// GetProvinceWards handles GET /api/v1/provinces/:code/wards
func (h *APIHandler) GetProvinceWards(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	provinceCode := c.Param("code")
	if provinceCode == "" {
		h.respondWithError(c, http.StatusBadRequest, "Province code is required")
		return
	}

	// Check if province exists
	_, err := h.dataService.GetProvince(provinceCode)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, "Province not found")
		return
	}

	search, typeFilter, limit, offset := h.parseQueryParams(c)

	wards, total := h.dataService.SearchWards(search, typeFilter, provinceCode, limit, offset)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    wards,
		Pagination: models.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Pages:  int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// Ward Handlers

// GetWards handles GET /api/v1/wards
func (h *APIHandler) GetWards(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	search, typeFilter, limit, offset := h.parseQueryParams(c)
	provinceCode := strings.TrimSpace(c.Query("province_code"))

	wards, total := h.dataService.SearchWards(search, typeFilter, provinceCode, limit, offset)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    wards,
		Pagination: models.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Pages:  int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// GetWard handles GET /api/v1/wards/:code
func (h *APIHandler) GetWard(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	code := c.Param("code")
	if code == "" {
		h.respondWithError(c, http.StatusBadRequest, "Ward code is required")
		return
	}

	ward, province, err := h.dataService.GetWardWithProvince(code)
	if err != nil {
		h.respondWithError(c, http.StatusNotFound, "Ward not found")
		return
	}

	// Create enhanced ward response with province info
	response := map[string]interface{}{
		"code":           ward.Code,
		"name":           ward.Name,
		"slug":           ward.Slug,
		"type":           ward.Type,
		"name_with_type": ward.NameWithType,
		"path":           ward.Path,
		"path_with_type": ward.PathWithType,
		"parent_code":    ward.ParentCode,
	}

	if province != nil {
		response["province"] = map[string]interface{}{
			"code":           province.Code,
			"name":           province.Name,
			"name_with_type": province.NameWithType,
			"type":           province.Type,
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// Search Handlers

// GlobalSearch handles GET /api/v1/search
func (h *APIHandler) GlobalSearch(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 2 {
		h.respondWithError(c, http.StatusBadRequest, "Search query must be at least 2 characters")
		return
	}

	entity := strings.TrimSpace(c.Query("entity"))
	if entity == "" {
		entity = "all"
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	results := h.dataService.GlobalSearch(query, entity, limit)

	c.JSON(http.StatusOK, models.SearchResponse{
		Success: true,
		Data:    results,
		Query:   query,
	})
}

// Utility Handlers

// ValidateAddress handles POST /api/v1/address/validate
func (h *APIHandler) ValidateAddress(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	var req models.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ward, valid := h.dataService.ValidateAddress(req.ProvinceCode, req.WardCode)

	response := models.ValidationResponse{
		Success: true,
		Valid:   valid,
	}

	if valid && ward != nil {
		response.Data = ward
		response.Message = "Address is valid"
	} else {
		response.Message = "Invalid address combination"
	}

	c.JSON(http.StatusOK, response)
}

// Health handles GET /api/v1/health
func (h *APIHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	status := "healthy"
	services := models.Services{
		DataLoader: "healthy",
	}

	if !h.dataService.IsDataLoaded() {
		status = "degraded"
		services.DataLoader = "not_loaded"
	}

	response := models.HealthResponse{
		Success:   true,
		Status:    status,
		Timestamp: time.Now(),
		Services:  services,
		Version:   h.version,
		Uptime:    uptime.String(),
	}

	statusCode := http.StatusOK
	if status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Stats handles GET /api/v1/stats (Admin endpoint)
func (h *APIHandler) Stats(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	stats := h.dataService.GetDataStats()
	stats["uptime"] = time.Since(h.startTime).String()
	stats["version"] = h.version

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// ReloadData handles POST /api/v1/admin/reload (Admin endpoint)
func (h *APIHandler) ReloadData(c *gin.Context) {
	err := h.dataService.ReloadData()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to reload data: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Data reloaded successfully",
		Data: map[string]interface{}{
			"reload_time": time.Now(),
			"stats":       h.dataService.GetDataStats(),
		},
	})
}

// GetProvinceTypes handles GET /api/v1/provinces/types
func (h *APIHandler) GetProvinceTypes(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	types := h.dataService.GetProvinceTypes()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    types,
	})
}

// GetWardTypes handles GET /api/v1/wards/types
func (h *APIHandler) GetWardTypes(c *gin.Context) {
	if !h.checkDataLoaded(c) {
		return
	}

	types := h.dataService.GetWardTypes()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    types,
	})
}

// NotFound handles 404 errors
func (h *APIHandler) NotFound(c *gin.Context) {
	h.respondWithError(c, http.StatusNotFound, "Endpoint not found")
}

// CORS preflight handler
func (h *APIHandler) Options(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Status(http.StatusOK)
}
