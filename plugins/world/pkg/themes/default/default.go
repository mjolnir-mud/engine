package _default

import "github.com/charmbracelet/lipgloss"

type theme struct{}

func (t theme) Name() string {
	return "default"
}

func (t theme) DefaultStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0c0c0"))
}

func (t theme) GetStyleFor(name string) lipgloss.Style {
	switch name {
	case "says":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF"))
	case "room_title":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	default:
		return t.DefaultStyle()
	}
}

var Theme = theme{}
