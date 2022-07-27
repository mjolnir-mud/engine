package theme

import "github.com/charmbracelet/lipgloss"

type Theme interface {
	Name() string
	DefaultStyle() lipgloss.Style
	GetStyleFor(name string) lipgloss.Style
}
