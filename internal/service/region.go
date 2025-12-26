package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	resourcecenter "github.com/alibabacloud-go/resourcecenter-20221201/client"
	"github.com/alibabacloud-go/tea/tea"

	"aliyun-tui-viewer/internal/i18n"
)

// RegionCache represents the cached region data for a profile
type RegionCache struct {
	Profile   string    `json:"profile"`
	Regions   []string  `json:"regions"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegionCacheFile represents the cache file structure
type RegionCacheFile struct {
	Profiles map[string]RegionCache `json:"profiles"`
}

// RegionService provides region-related operations
type RegionService struct {
	accessKeyID     string
	accessKeySecret string
	profile         string
	cacheExpiry     time.Duration
}

// NewRegionService creates a new RegionService
func NewRegionService(accessKeyID, accessKeySecret, profile string) *RegionService {
	return &RegionService{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		profile:         profile,
		cacheExpiry:     7 * 24 * time.Hour, // 7 days
	}
}

// GetRegionsWithResources returns the list of regions where the account has resources
func (s *RegionService) GetRegionsWithResources() ([]string, error) {
	// Check cache first
	if cached := s.loadCache(); cached != nil && !s.isCacheExpired(cached) {
		return cached.Regions, nil
	}

	// Fetch from API
	regions, err := s.fetchRegionsFromAPI()
	if err != nil {
		return nil, err
	}

	// Save to cache
	s.saveCache(regions)

	return regions, nil
}

// ForceRefresh forces a refresh of the region list, bypassing cache
func (s *RegionService) ForceRefresh() ([]string, error) {
	regions, err := s.fetchRegionsFromAPI()
	if err != nil {
		return nil, err
	}

	s.saveCache(regions)
	return regions, nil
}

// fetchRegionsFromAPI fetches regions from Aliyun Resource Center API
func (s *RegionService) fetchRegionsFromAPI() ([]string, error) {
	// Create Resource Center client
	config := &openapi.Config{
		AccessKeyId:     tea.String(s.accessKeyID),
		AccessKeySecret: tea.String(s.accessKeySecret),
		Endpoint:        tea.String("resourcecenter.aliyuncs.com"),
	}

	client, err := resourcecenter.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("creating resource center client: %w", err)
	}

	regionSet := make(map[string]bool)
	var nextToken *string

	for {
		request := &resourcecenter.SearchResourcesRequest{
			MaxResults: tea.Int32(500),
			NextToken:  nextToken,
		}

		response, err := client.SearchResources(request)
		if err != nil {
			return nil, fmt.Errorf("searching resources: %w", err)
		}

		if response.Body == nil || response.Body.Resources == nil {
			break
		}

		for _, resource := range response.Body.Resources {
			if resource.RegionId != nil && *resource.RegionId != "" {
				regionSet[*resource.RegionId] = true
			}
		}

		// Check for next page
		if response.Body.NextToken == nil || *response.Body.NextToken == "" {
			break
		}
		nextToken = response.Body.NextToken
	}

	// Convert set to sorted slice
	regions := make([]string, 0, len(regionSet))
	for region := range regionSet {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	return regions, nil
}

// getCachePath returns the path to the cache file
func (s *RegionService) getCachePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("getting current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".aliyun", "region_cache.json"), nil
}

// loadCache loads the cache for the current profile
func (s *RegionService) loadCache() *RegionCache {
	cachePath, err := s.getCachePath()
	if err != nil {
		return nil
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil
	}

	var cacheFile RegionCacheFile
	if err := json.Unmarshal(data, &cacheFile); err != nil {
		return nil
	}

	if cache, ok := cacheFile.Profiles[s.profile]; ok {
		return &cache
	}

	return nil
}

// saveCache saves the regions to cache for the current profile
func (s *RegionService) saveCache(regions []string) error {
	cachePath, err := s.getCachePath()
	if err != nil {
		return err
	}

	// Load existing cache file or create new one
	var cacheFile RegionCacheFile

	data, err := os.ReadFile(cachePath)
	if err == nil {
		json.Unmarshal(data, &cacheFile)
	}

	if cacheFile.Profiles == nil {
		cacheFile.Profiles = make(map[string]RegionCache)
	}

	// Update cache for current profile
	cacheFile.Profiles[s.profile] = RegionCache{
		Profile:   s.profile,
		Regions:   regions,
		UpdatedAt: time.Now(),
	}

	// Write cache file
	updatedData, err := json.MarshalIndent(cacheFile, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling cache: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	if err := os.WriteFile(cachePath, updatedData, 0644); err != nil {
		return fmt.Errorf("writing cache file: %w", err)
	}

	return nil
}

// isCacheExpired checks if the cache is expired
func (s *RegionService) isCacheExpired(cache *RegionCache) bool {
	return time.Since(cache.UpdatedAt) > s.cacheExpiry
}

// ClearCache clears the cache for the current profile
func (s *RegionService) ClearCache() error {
	cachePath, err := s.getCachePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil // No cache file, nothing to clear
	}

	var cacheFile RegionCacheFile
	if err := json.Unmarshal(data, &cacheFile); err != nil {
		return nil
	}

	delete(cacheFile.Profiles, s.profile)

	updatedData, err := json.MarshalIndent(cacheFile, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, updatedData, 0644)
}

// GetRegionDisplayName returns a human-readable name for a region ID
func GetRegionDisplayName(regionID string) string {
	return i18n.GetRegionDisplayName(regionID)
}

