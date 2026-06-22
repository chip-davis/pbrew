package ui

import "charm.land/lipgloss/v2"

// Palette anchored on the Homebrew orange (#FBB040). Tweak as needed.

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FBB040"))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8BE9FD")).
			MarginTop(1)

	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C"))

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5555"))
)
