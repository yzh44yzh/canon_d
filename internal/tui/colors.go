package tui

import (
	"image/color"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

// ColorScheme assigns terminal colors to the UI's five semantic roles.
// Ordinary text and backgrounds inherit the terminal's defaults.
type ColorScheme struct {
	Border    color.Color
	Active    color.Color
	Muted     color.Color
	Correct   color.Color
	Incorrect color.Color
}

// DefaultColorScheme uses the terminal's palette rather than fixed RGB colors.
var DefaultColorScheme = ColorScheme{
	Border:    lipgloss.Color("8"),
	Active:    lipgloss.Color("6"),
	Muted:     lipgloss.Color("8"),
	Correct:   lipgloss.Color("2"),
	Incorrect: lipgloss.Color("1"),
}

type uiStyles struct {
	panel     lipgloss.Style
	active    lipgloss.Style
	muted     lipgloss.Style
	correct   lipgloss.Style
	incorrect lipgloss.Style
	input     textinput.Styles
}

func (scheme ColorScheme) styles() uiStyles {
	muted := lipgloss.NewStyle().Foreground(scheme.Muted)
	input := textinput.DefaultStyles(false)
	state := textinput.StyleState{
		Text:        lipgloss.NewStyle(),
		Placeholder: muted,
		Suggestion:  muted,
		Prompt:      lipgloss.NewStyle().Foreground(scheme.Active),
	}

	input.Focused, input.Blurred = state, state
	input.Cursor.Color = scheme.Active

	return uiStyles{
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(scheme.Border),
		active:    lipgloss.NewStyle().Foreground(scheme.Active).Bold(true).Reverse(true),
		muted:     muted,
		correct:   lipgloss.NewStyle().Foreground(scheme.Correct),
		incorrect: lipgloss.NewStyle().Foreground(scheme.Incorrect).Underline(true),
		input:     input,
	}
}
