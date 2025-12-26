package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#06B6D4") // Cyan
	AccentColor    = lipgloss.Color("#F59E0B") // Amber

	// Status colors
	SuccessColor = lipgloss.Color("#10B981") // Green
	WarningColor = lipgloss.Color("#F59E0B") // Amber
	ErrorColor   = lipgloss.Color("#EF4444") // Red
	InfoColor    = lipgloss.Color("#3B82F6") // Blue

	// Neutral colors
	TextColor       = lipgloss.Color("#E5E7EB") // Light gray
	SubtleTextColor = lipgloss.Color("#9CA3AF") // Gray
	MutedTextColor  = lipgloss.Color("#6B7280") // Dark gray
	BorderColor     = lipgloss.Color("#374151") // Dark border
	HighlightBg     = lipgloss.Color("#1F2937") // Highlight background

	// Special colors
	SelectedBg     = lipgloss.Color("#374151")
	SearchMatchBg  = lipgloss.Color("#854D0E") // Dark yellow
	CurrentMatchBg = lipgloss.Color("#CA8A04") // Bright yellow
)

// Styles contains all application styles
type Styles struct {
	// App container
	App lipgloss.Style

	// Header and title
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Description lipgloss.Style

	// Table styles
	TableHeader      lipgloss.Style
	TableCell        lipgloss.Style
	TableSelectedRow lipgloss.Style
	TableBorder      lipgloss.Style

	// List/Menu styles
	MenuItem         lipgloss.Style
	MenuItemSelected lipgloss.Style
	MenuShortcut     lipgloss.Style

	// Detail view styles
	JSONKey     lipgloss.Style
	JSONValue   lipgloss.Style
	JSONString  lipgloss.Style
	JSONNumber  lipgloss.Style
	JSONBoolean lipgloss.Style
	JSONNull    lipgloss.Style

	// Search styles
	SearchBar      lipgloss.Style
	SearchMatch    lipgloss.Style
	SearchCurrent  lipgloss.Style
	SearchLabel    lipgloss.Style
	SearchNoResult lipgloss.Style

	// Mode line styles
	ModeLine        lipgloss.Style
	ModeLineProfile lipgloss.Style
	ModeLineHelp    lipgloss.Style
	ModeLineInfo    lipgloss.Style

	// Modal styles
	ModalOverlay lipgloss.Style
	ModalContent lipgloss.Style
	ModalTitle   lipgloss.Style
	ModalButton  lipgloss.Style

	// Status styles
	StatusSuccess lipgloss.Style
	StatusWarning lipgloss.Style
	StatusError   lipgloss.Style
	StatusInfo    lipgloss.Style

	// Spinner/Loading
	Spinner lipgloss.Style

	// Help
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style
	HelpSep  lipgloss.Style

	// Border
	Border        lipgloss.Style
	FocusedBorder lipgloss.Style
}

// DefaultStyles returns the default style set
func DefaultStyles() *Styles {
	s := &Styles{}

	// App container
	s.App = lipgloss.NewStyle()

	// Header and title
	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(SecondaryColor)

	s.Description = lipgloss.NewStyle().
		Foreground(SubtleTextColor)

	// Table styles
	s.TableHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(BorderColor)

	s.TableCell = lipgloss.NewStyle().
		Foreground(TextColor).
		Padding(0, 1)

	s.TableSelectedRow = lipgloss.NewStyle().
		Background(SelectedBg).
		Foreground(TextColor).
		Bold(true)

	s.TableBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	// List/Menu styles
	s.MenuItem = lipgloss.NewStyle().
		Foreground(TextColor).
		PaddingLeft(2)

	s.MenuItemSelected = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		PaddingLeft(2)

	s.MenuShortcut = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	// Detail view styles
	s.JSONKey = lipgloss.NewStyle().
		Foreground(SecondaryColor)

	s.JSONValue = lipgloss.NewStyle().
		Foreground(TextColor)

	s.JSONString = lipgloss.NewStyle().
		Foreground(SuccessColor)

	s.JSONNumber = lipgloss.NewStyle().
		Foreground(AccentColor)

	s.JSONBoolean = lipgloss.NewStyle().
		Foreground(PrimaryColor)

	s.JSONNull = lipgloss.NewStyle().
		Foreground(MutedTextColor).
		Italic(true)

	// Search styles
	s.SearchBar = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 1)

	s.SearchMatch = lipgloss.NewStyle().
		Background(SearchMatchBg).
		Foreground(TextColor)

	s.SearchCurrent = lipgloss.NewStyle().
		Background(CurrentMatchBg).
		Foreground(lipgloss.Color("#000000")).
		Bold(true)

	s.SearchLabel = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	s.SearchNoResult = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Italic(true)

	// Mode line styles
	s.ModeLine = lipgloss.NewStyle().
		Background(HighlightBg).
		Foreground(TextColor).
		Padding(0, 1)

	s.ModeLineProfile = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	s.ModeLineHelp = lipgloss.NewStyle().
		Foreground(SubtleTextColor)

	s.ModeLineInfo = lipgloss.NewStyle().
		Foreground(SecondaryColor)

	// Modal styles
	s.ModalOverlay = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000"))

	s.ModalContent = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2).
		Background(HighlightBg)

	s.ModalTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	s.ModalButton = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(BorderColor).
		Padding(0, 2).
		MarginRight(1)

	// Status styles
	s.StatusSuccess = lipgloss.NewStyle().
		Foreground(SuccessColor)

	s.StatusWarning = lipgloss.NewStyle().
		Foreground(WarningColor)

	s.StatusError = lipgloss.NewStyle().
		Foreground(ErrorColor)

	s.StatusInfo = lipgloss.NewStyle().
		Foreground(InfoColor)

	// Spinner/Loading
	s.Spinner = lipgloss.NewStyle().
		Foreground(PrimaryColor)

	// Help
	s.HelpKey = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	s.HelpDesc = lipgloss.NewStyle().
		Foreground(SubtleTextColor)

	s.HelpSep = lipgloss.NewStyle().
		Foreground(MutedTextColor)

	// Border
	s.Border = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	s.FocusedBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor)

	return s
}

// GlobalStyles is the default style instance
var GlobalStyles = DefaultStyles()

// Helper functions for common styling operations

// RenderTitle renders a title with the default style
func RenderTitle(title string) string {
	return GlobalStyles.Title.Render(title)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return GlobalStyles.StatusError.Render("Error: " + msg)
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return GlobalStyles.StatusSuccess.Render(msg)
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return GlobalStyles.StatusInfo.Render(msg)
}

// CenterHorizontally centers content horizontally within the given width
func CenterHorizontally(content string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, content)
}

// CenterVertically centers content vertically within the given height
func CenterVertically(content string, height int) string {
	return lipgloss.PlaceVertical(height, lipgloss.Center, content)
}

// Center centers content both horizontally and vertically
func Center(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

