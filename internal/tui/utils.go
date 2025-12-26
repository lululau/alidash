package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/atotto/clipboard"

	"aliyun-tui-viewer/internal/config"
)

// CopyToClipboard copies data to clipboard
func CopyToClipboard(data interface{}) tea.Cmd {
	return func() tea.Msg {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("marshaling data: %w", err)}
		}

		if err := clipboard.WriteAll(string(jsonData)); err != nil {
			return ErrorMsg{Err: fmt.Errorf("copying to clipboard: %w", err)}
		}

		return CopiedMsg{}
	}
}

// OpenInEditor opens data in the configured external editor
func OpenInEditor(data interface{}) tea.Cmd {
	return tea.ExecProcess(createEditorCmd(data), func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("editor closed with error: %w", err)}
		}
		return EditorClosedMsg{}
	})
}

// createEditorCmd creates the editor command
func createEditorCmd(data interface{}) *exec.Cmd {
	// Get editor from config or environment
	editor, err := config.GetEditor()
	if err != nil || editor == "" {
		editor = "nvim" // Default to nvim
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "tali-*.json")
	if err != nil {
		return exec.Command("echo", "Error creating temp file")
	}

	// Write data to temp file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return exec.Command("echo", "Error marshaling data")
	}

	if _, err := tmpFile.Write(jsonData); err != nil {
		return exec.Command("echo", "Error writing to temp file")
	}
	tmpFile.Close()

	return exec.Command(editor, tmpFile.Name())
}

// OpenInPager opens data in the configured external pager
func OpenInPager(data interface{}) tea.Cmd {
	return tea.ExecProcess(createPagerCmd(data), func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("pager closed with error: %w", err)}
		}
		return EditorClosedMsg{} // Reuse EditorClosedMsg for pager as well
	})
}

// createPagerCmd creates the pager command
func createPagerCmd(data interface{}) *exec.Cmd {
	// Get pager from config or environment
	pager, err := config.GetPager()
	if err != nil || pager == "" {
		pager = "less" // Default to less
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "tali-*.json")
	if err != nil {
		return exec.Command("echo", "Error creating temp file")
	}

	// Write data to temp file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return exec.Command("echo", "Error marshaling data")
	}

	if _, err := tmpFile.Write(jsonData); err != nil {
		return exec.Command("echo", "Error writing to temp file")
	}
	tmpFile.Close()

	return exec.Command(pager, tmpFile.Name())
}

// SwitchProfile switches to a different profile and reinitializes clients
func SwitchProfile(profileName string) tea.Cmd {
	return func() tea.Msg {
		// Switch profile in config
		if err := config.SwitchProfile(profileName); err != nil {
			return ErrorMsg{Err: fmt.Errorf("switching profile: %w", err)}
		}

		return ProfileSwitchedMsg{ProfileName: profileName}
	}
}

// LoadProfileList loads the list of available profiles
func LoadProfileList() tea.Cmd {
	return func() tea.Msg {
		profiles, err := config.ListAllProfiles()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("loading profiles: %w", err)}
		}

		currentProfile, err := config.GetCurrentProfileName()
		if err != nil {
			currentProfile = "default"
		}

		return ProfileListLoadedMsg{
			Profiles:       profiles,
			CurrentProfile: currentProfile,
		}
	}
}

