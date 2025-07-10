package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"vietnam-admin-api/models"
)

// DataService handles loading and accessing Vietnamese administrative data
type DataService struct {
	provinces models.ProvinceData
	wards     models.WardData
	mu        sync.RWMutex
	loadTime  time.Time
	dataPath  string
}

// NewDataService creates a new DataService instance
func NewDataService(dataPath string) *DataService {
	return &DataService{
		dataPath: dataPath,
	}
}

// LoadData loads JSON data from files into memory
func (ds *DataService) LoadData() error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	log.Println("Loading Vietnamese administrative data...")
	startTime := time.Now()

	// Load provinces
	provinceFile := filepath.Join(ds.dataPath, "province.json")
	provinceData, err := os.ReadFile(provinceFile)
	if err != nil {
		return fmt.Errorf("failed to read province.json: %w", err)
	}

	provinces, err := models.UnmarshalProvinceData(provinceData)
	if err != nil {
		return fmt.Errorf("failed to parse province.json: %w", err)
	}

	// Load wards
	wardFile := filepath.Join(ds.dataPath, "ward.json")
	wardData, err := os.ReadFile(wardFile)
	if err != nil {
		return fmt.Errorf("failed to read ward.json: %w", err)
	}

	wards, err := models.UnmarshalWardData(wardData)
	if err != nil {
		return fmt.Errorf("failed to parse ward.json: %w", err)
	}

	// Set loaded data
	ds.provinces = provinces
	ds.wards = wards
	ds.loadTime = time.Now()

	loadDuration := time.Since(startTime)
	log.Printf("Data loaded successfully in %v - Provinces: %d, Wards: %d",
		loadDuration, len(ds.provinces), len(ds.wards))

	return nil
}

// ReloadData reloads data from JSON files
func (ds *DataService) ReloadData() error {
	log.Println("Reloading data...")
	return ds.LoadData()
}

// GetLoadTime returns when data was last loaded
func (ds *DataService) GetLoadTime() time.Time {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.loadTime
}

// IsDataLoaded checks if data has been loaded
func (ds *DataService) IsDataLoaded() bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return len(ds.provinces) > 0 && len(ds.wards) > 0
}

// GetDataStats returns statistics about loaded data
func (ds *DataService) GetDataStats() map[string]interface{} {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	stats := map[string]interface{}{
		"provinces": len(ds.provinces),
		"wards":     len(ds.wards),
		"load_time": ds.loadTime,
		"is_loaded": ds.IsDataLoaded(),
	}

	if ds.IsDataLoaded() {
		// Count by province types
		provinceTypes := make(map[string]int)
		for _, province := range ds.provinces {
			provinceTypes[province.Type]++
		}

		// Count by ward types
		wardTypes := make(map[string]int)
		for _, ward := range ds.wards {
			wardTypes[ward.Type]++
		}

		stats["province_types"] = provinceTypes
		stats["ward_types"] = wardTypes
	}

	return stats
}

// Province Methods

// GetAllProvinces returns all provinces
func (ds *DataService) GetAllProvinces() []models.Province {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	provinces := ds.provinces.ToSlice()
	sort.Slice(provinces, func(i, j int) bool {
		return provinces[i].Name < provinces[j].Name
	})
	return provinces
}

// GetProvince returns a province by code
func (ds *DataService) GetProvince(code string) (*models.Province, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	province, exists := ds.provinces[code]
	if !exists {
		return nil, fmt.Errorf("province with code %s not found", code)
	}
	return &province, nil
}

// SearchProvinces searches provinces with filters and pagination
func (ds *DataService) SearchProvinces(search, typeFilter string, limit, offset int) ([]models.Province, int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Filter provinces
	filteredProvinces := ds.provinces.ToSliceWithFilters(search, typeFilter)

	// Sort by name
	sort.Slice(filteredProvinces, func(i, j int) bool {
		return filteredProvinces[i].Name < filteredProvinces[j].Name
	})

	// Apply pagination
	paginatedData, total := models.PaginateSlice(filteredProvinces, limit, offset)
	return paginatedData.([]models.Province), total
}

// GetProvinceTypes returns all unique province types
func (ds *DataService) GetProvinceTypes() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	typeMap := make(map[string]bool)
	for _, province := range ds.provinces {
		typeMap[province.Type] = true
	}

	types := make([]string, 0, len(typeMap))
	for t := range typeMap {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// Ward Methods

// GetAllWards returns all wards
func (ds *DataService) GetAllWards() []models.Ward {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	wards := ds.wards.ToSlice()
	sort.Slice(wards, func(i, j int) bool {
		return wards[i].Name < wards[j].Name
	})
	return wards
}

// GetWard returns a ward by code
func (ds *DataService) GetWard(code string) (*models.Ward, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	ward, exists := ds.wards[code]
	if !exists {
		return nil, fmt.Errorf("ward with code %s not found", code)
	}
	return &ward, nil
}

// GetWardsByProvince returns wards belonging to a specific province
func (ds *DataService) GetWardsByProvince(provinceCode string) []models.Ward {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	wards := ds.wards.ToSliceWithFilters("", "", provinceCode)
	sort.Slice(wards, func(i, j int) bool {
		return wards[i].Name < wards[j].Name
	})
	return wards
}

// SearchWards searches wards with filters and pagination
func (ds *DataService) SearchWards(search, typeFilter, provinceCode string, limit, offset int) ([]models.Ward, int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Filter wards
	filteredWards := ds.wards.ToSliceWithFilters(search, typeFilter, provinceCode)

	// Sort by name
	sort.Slice(filteredWards, func(i, j int) bool {
		return filteredWards[i].Name < filteredWards[j].Name
	})

	// Apply pagination
	paginatedData, total := models.PaginateSlice(filteredWards, limit, offset)
	return paginatedData.([]models.Ward), total
}

// GetWardTypes returns all unique ward types
func (ds *DataService) GetWardTypes() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	typeMap := make(map[string]bool)
	for _, ward := range ds.wards {
		typeMap[ward.Type] = true
	}

	types := make([]string, 0, len(typeMap))
	for t := range typeMap {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// Search Methods

// GlobalSearch performs a global search across provinces and wards
func (ds *DataService) GlobalSearch(query string, entity string, limit int) models.SearchData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	result := models.SearchData{
		Provinces: []models.Province{},
		Wards:     []models.Ward{},
	}

	if entity == "all" || entity == "province" {
		provinces := ds.provinces.ToSliceWithFilters(query, "")
		sort.Slice(provinces, func(i, j int) bool {
			return provinces[i].Name < provinces[j].Name
		})
		if len(provinces) > limit {
			provinces = provinces[:limit]
		}
		result.Provinces = provinces
	}

	if entity == "all" || entity == "ward" {
		wards := ds.wards.ToSliceWithFilters(query, "", "")
		sort.Slice(wards, func(i, j int) bool {
			return wards[i].Name < wards[j].Name
		})
		if len(wards) > limit {
			wards = wards[:limit]
		}
		result.Wards = wards
	}

	return result
}

// ValidateAddress validates if a ward belongs to a province
func (ds *DataService) ValidateAddress(provinceCode, wardCode string) (*models.Ward, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Check if province exists
	_, provinceExists := ds.provinces[provinceCode]
	if !provinceExists {
		return nil, false
	}

	// Check if ward exists and belongs to the province
	ward, wardExists := ds.wards[wardCode]
	if !wardExists {
		return nil, false
	}

	if ward.ParentCode != provinceCode {
		return nil, false
	}

	return &ward, true
}

// GetWardWithProvince returns ward with province information
func (ds *DataService) GetWardWithProvince(wardCode string) (*models.Ward, *models.Province, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	ward, exists := ds.wards[wardCode]
	if !exists {
		return nil, nil, fmt.Errorf("ward with code %s not found", wardCode)
	}

	province, exists := ds.provinces[ward.ParentCode]
	if !exists {
		return &ward, nil, fmt.Errorf("province with code %s not found", ward.ParentCode)
	}

	return &ward, &province, nil
}
