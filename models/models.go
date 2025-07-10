package models

import (
	"encoding/json"
	"strings"
	"time"
)

// Province represents a Vietnamese province/city
type Province struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	NameWithType string `json:"name_with_type"`
}

// Ward represents a Vietnamese ward/commune/town
type Ward struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	NameWithType string `json:"name_with_type"`
	Path         string `json:"path"`
	PathWithType string `json:"path_with_type"`
	ParentCode   string `json:"parent_code"`
}

// ProvinceData represents the structure of province.json
type ProvinceData map[string]Province

// WardData represents the structure of ward.json
type WardData map[string]Ward

// Response structures for API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Message    string      `json:"message,omitempty"`
}

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Pages  int `json:"pages"`
}

type SearchResponse struct {
	Success bool       `json:"success"`
	Data    SearchData `json:"data"`
	Query   string     `json:"query"`
	Message string     `json:"message,omitempty"`
}

type SearchData struct {
	Provinces []Province `json:"provinces"`
	Wards     []Ward     `json:"wards"`
}

type ValidationRequest struct {
	ProvinceCode string `json:"province_code" binding:"required"`
	WardCode     string `json:"ward_code" binding:"required"`
}

type ValidationResponse struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid"`
	Data    *Ward  `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Success   bool      `json:"success"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Services  Services  `json:"services"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

type Services struct {
	DataLoader string `json:"data_loader"`
	Cache      string `json:"cache,omitempty"`
}

// Search methods for Province
func (p Province) MatchesQuery(query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(p.Name), query) ||
		strings.Contains(strings.ToLower(p.Slug), query) ||
		strings.Contains(strings.ToLower(p.NameWithType), query)
}

// Search methods for Ward
func (w Ward) MatchesQuery(query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(w.Name), query) ||
		strings.Contains(strings.ToLower(w.Slug), query) ||
		strings.Contains(strings.ToLower(w.NameWithType), query) ||
		strings.Contains(strings.ToLower(w.Path), query) ||
		strings.Contains(strings.ToLower(w.PathWithType), query)
}

// Filter methods for Province
func (p Province) MatchesType(typeFilter string) bool {
	return typeFilter == "" || p.Type == typeFilter
}

// Filter methods for Ward
func (w Ward) MatchesType(typeFilter string) bool {
	return typeFilter == "" || w.Type == typeFilter
}

func (w Ward) MatchesParentCode(parentCode string) bool {
	return parentCode == "" || w.ParentCode == parentCode
}

// Convert map to slice methods
func (pd ProvinceData) ToSlice() []Province {
	provinces := make([]Province, 0, len(pd))
	for _, province := range pd {
		provinces = append(provinces, province)
	}
	return provinces
}

func (wd WardData) ToSlice() []Ward {
	wards := make([]Ward, 0, len(wd))
	for _, ward := range wd {
		wards = append(wards, ward)
	}
	return wards
}

// Convert map to slice with filters
func (pd ProvinceData) ToSliceWithFilters(search, typeFilter string) []Province {
	provinces := []Province{}
	for _, province := range pd {
		if (search == "" || province.MatchesQuery(search)) &&
			province.MatchesType(typeFilter) {
			provinces = append(provinces, province)
		}
	}
	return provinces
}

func (wd WardData) ToSliceWithFilters(search, typeFilter, parentCode string) []Ward {
	wards := []Ward{}
	for _, ward := range wd {
		if (search == "" || ward.MatchesQuery(search)) &&
			ward.MatchesType(typeFilter) &&
			ward.MatchesParentCode(parentCode) {
			wards = append(wards, ward)
		}
	}
	return wards
}

// Pagination helper
func PaginateSlice(slice interface{}, limit, offset int) (interface{}, int) {
	switch s := slice.(type) {
	case []Province:
		total := len(s)
		if offset >= total {
			return []Province{}, total
		}
		end := offset + limit
		if end > total {
			end = total
		}
		return s[offset:end], total
	case []Ward:
		total := len(s)
		if offset >= total {
			return []Ward{}, total
		}
		end := offset + limit
		if end > total {
			end = total
		}
		return s[offset:end], total
	default:
		return slice, 0
	}
}

// JSON unmarshal helpers
func UnmarshalProvinceData(data []byte) (ProvinceData, error) {
	var provinces ProvinceData
	err := json.Unmarshal(data, &provinces)
	return provinces, err
}

func UnmarshalWardData(data []byte) (WardData, error) {
	var wards WardData
	err := json.Unmarshal(data, &wards)
	return wards, err
}
