package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// ConfigProfile represents a single profile in the Aliyun CLI config
type ConfigProfile struct {
	Name            string `json:"name"`
	Mode            string `json:"mode"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	RegionID        string `json:"region_id"`
	OssEndpoint     string `json:"oss_endpoint,omitempty"` // Custom field for OSS endpoint
	// Other fields like output_format, language can be added if needed
}

// AliyunConfig represents the structure of ~/.aliyun/config.json
type AliyunConfig struct {
	Current  string          `json:"current"`
	Profiles []ConfigProfile `json:"profiles"`
	Editor   string          `json:"editor,omitempty"` // Global editor command
	Pager    string          `json:"pager,omitempty"`  // Global pager command
}

// Config holds the application configuration
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	OssEndpoint     string
	Editor          string
	Pager           string
}

// LoadAliyunConfig loads configuration from ~/.aliyun/config.json
func LoadAliyunConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".aliyun", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read aliyun config file at %s: %w", configPath, err)
	}

	var config AliyunConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aliyun config file %s: %w", configPath, err)
	}

	if len(config.Profiles) == 0 {
		return nil, fmt.Errorf("no profiles found in aliyun config file: %s", configPath)
	}

	activeProfileName := config.Current
	if activeProfileName == "" {
		if len(config.Profiles) == 1 {
			activeProfileName = config.Profiles[0].Name
		} else {
			for _, p := range config.Profiles {
				if p.Name == "default" {
					activeProfileName = "default"
					break
				}
			}
			if activeProfileName == "" { // If no "default" profile and current is not set
				return nil, fmt.Errorf("no current profile specified in %s, and no 'default' profile found. Please specify a current profile or name one 'default'", configPath)
			}
		}
	}

	var activeProfile *ConfigProfile
	for i := range config.Profiles {
		if config.Profiles[i].Name == activeProfileName {
			activeProfile = &config.Profiles[i]
			break
		}
	}

	if activeProfile == nil {
		return nil, fmt.Errorf("current profile '%s' not found in aliyun config file: %s", activeProfileName, configPath)
	}

	if activeProfile.AccessKeyID == "" || activeProfile.AccessKeySecret == "" || activeProfile.RegionID == "" {
		return nil, fmt.Errorf("profile '%s' in %s is missing access_key_id, access_key_secret, or region_id", activeProfile.Name, configPath)
	}

	// Resolve OSS Endpoint
	ossEndpoint := activeProfile.OssEndpoint
	if ossEndpoint == "" && activeProfile.RegionID != "" {
		ossEndpoint = fmt.Sprintf("oss-%s.aliyuncs.com", activeProfile.RegionID)
	}

	if ossEndpoint == "" {
		return nil, fmt.Errorf("OSS endpoint could not be determined for profile '%s'. "+
			"Please either set 'oss_endpoint' in your profile in %s, "+
			"or ensure 'region_id' is set for default construction (e.g., oss-<region_id>.aliyuncs.com)",
			activeProfile.Name, configPath)
	}

	return &Config{
		AccessKeyID:     activeProfile.AccessKeyID,
		AccessKeySecret: activeProfile.AccessKeySecret,
		RegionID:        activeProfile.RegionID,
		OssEndpoint:     ossEndpoint,
		Editor:          config.Editor,
		Pager:           config.Pager,
	}, nil
}

// GetCurrentProfileName returns the name of the current active profile
func GetCurrentProfileName() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".aliyun", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read aliyun config file at %s: %w", configPath, err)
	}

	var config AliyunConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return "", fmt.Errorf("failed to parse aliyun config file %s: %w", configPath, err)
	}

	if config.Current != "" {
		return config.Current, nil
	}

	// Fallback logic
	if len(config.Profiles) == 1 {
		return config.Profiles[0].Name, nil
	}

	for _, p := range config.Profiles {
		if p.Name == "default" {
			return "default", nil
		}
	}

	if len(config.Profiles) > 0 {
		return config.Profiles[0].Name, nil
	}

	return "", fmt.Errorf("no profiles found")
}

// ListAllProfiles returns all available profile names
func ListAllProfiles() ([]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".aliyun", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read aliyun config file at %s: %w", configPath, err)
	}

	var config AliyunConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aliyun config file %s: %w", configPath, err)
	}

	var profiles []string
	for _, profile := range config.Profiles {
		profiles = append(profiles, profile.Name)
	}

	return profiles, nil
}

// SwitchProfile switches to the specified profile
func SwitchProfile(profileName string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".aliyun", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read aliyun config file at %s: %w", configPath, err)
	}

	var config AliyunConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to parse aliyun config file %s: %w", configPath, err)
	}

	// Check if profile exists
	profileExists := false
	for _, profile := range config.Profiles {
		if profile.Name == profileName {
			profileExists = true
			break
		}
	}

	if !profileExists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Update current profile
	config.Current = profileName

	// Write back to file
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetEditor returns the editor command to use, following the priority:
// 1. Config file "editor" field
// 2. VISUAL environment variable
// 3. EDITOR environment variable
// 4. Default to "vim"
func GetEditor() (string, error) {
	config, err := LoadAliyunConfig()
	if err != nil {
		return "", err
	}

	// First check config file
	if config.Editor != "" {
		return config.Editor, nil
	}

	// Then check VISUAL environment variable
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual, nil
	}

	// Then check EDITOR environment variable
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Default to vim
	return "vim", nil
}

// GetPager returns the pager command to use, following the priority:
// 1. Config file "pager" field
// 2. PAGER environment variable
// 3. Default to "less"
func GetPager() (string, error) {
	config, err := LoadAliyunConfig()
	if err != nil {
		return "", err
	}

	// First check config file
	if config.Pager != "" {
		return config.Pager, nil
	}

	// Then check PAGER environment variable
	if pager := os.Getenv("PAGER"); pager != "" {
		return pager, nil
	}

	// Default to less
	return "less", nil
}

// InputHistory manages input history for the finder
type InputHistory struct {
	Items    []string `json:"items"`
	MaxItems int      `json:"-"`
}

// historyFilePath returns the path to the history file
func historyFilePath() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, ".aliyun", "tali_history.json")
}

// LoadInputHistory loads input history from file
func LoadInputHistory() *InputHistory {
	h := &InputHistory{
		Items:    []string{},
		MaxItems: 100,
	}

	path := historyFilePath()
	if path == "" {
		return h
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return h
	}

	_ = json.Unmarshal(data, h)
	h.MaxItems = 100
	return h
}

// Save saves the input history to file
func (h *InputHistory) Save() error {
	path := historyFilePath()
	if path == "" {
		return fmt.Errorf("could not determine history file path")
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Add adds an item to the history
func (h *InputHistory) Add(item string) {
	if item == "" {
		return
	}

	// Remove duplicates
	newItems := []string{item}
	for _, existing := range h.Items {
		if existing != item {
			newItems = append(newItems, existing)
		}
	}

	// Limit to max items
	if len(newItems) > h.MaxItems {
		newItems = newItems[:h.MaxItems]
	}

	h.Items = newItems
}

// Get returns the item at the given index (0 = most recent)
func (h *InputHistory) Get(index int) string {
	if index < 0 || index >= len(h.Items) {
		return ""
	}
	return h.Items[index]
}

// Len returns the number of items in history
func (h *InputHistory) Len() int {
	return len(h.Items)
}
