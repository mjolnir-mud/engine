package default_theme

import "github.com/charmbracelet/lipgloss"

type defualtTheme struct{}

func (t *defualtTheme) Name() string {
	return "default"
}

func (t *defualtTheme) GetStyleFor(name string) lipgloss.Style {
	switch name {
	case "command":
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("white"))
	case "ordered-list-item":
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("white")).
			PaddingLeft(4)
	default:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("white"))
	}
}

var Theme = &defualtTheme{}
