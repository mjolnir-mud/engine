package theme

import "github.com/charmbracelet/lipgloss"

type Theme interface {
	Name() string
	GetStyleFor(name string) lipgloss.Style
}
