package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	Pink   = lipgloss.AdaptiveColor{Light: "#D63384", Dark: "#FF79C6"}
	Purple = lipgloss.AdaptiveColor{Light: "#6F42C1", Dark: "#BD93F9"}
	Blue   = lipgloss.AdaptiveColor{Light: "#0D6EFD", Dark: "#8BE9FD"}
	Green  = lipgloss.AdaptiveColor{Light: "#198754", Dark: "#50FA7B"}
	Yellow = lipgloss.AdaptiveColor{Light: "#B58100", Dark: "#F1FA8C"}
	Peach  = lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FFB86C"}
	Red    = lipgloss.AdaptiveColor{Light: "#DC3545", Dark: "#FF5555"}
	Text   = lipgloss.AdaptiveColor{Light: "#212529", Dark: "#F8F8F2"}
	Muted  = lipgloss.AdaptiveColor{Light: "#6C757D", Dark: "#6272A4"}
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Blue)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Pink).
			Padding(1, 2)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	PinkStyle = lipgloss.NewStyle().
			Foreground(Pink)

	PurpleStyle = lipgloss.NewStyle().
			Foreground(Purple)

	BlueStyle = lipgloss.NewStyle().
			Foreground(Blue)

	GreenStyle = lipgloss.NewStyle().
			Foreground(Green)

	PeachStyle = lipgloss.NewStyle().
			Foreground(Peach)

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(Purple).
				Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Muted)

	SectionStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(Pink).
			Bold(true)

	SoftValueStyle = lipgloss.NewStyle().
			Foreground(Text)
)

func GradientText(text, from, to string) string {
	start, err := colorful.Hex(from)
	if err != nil {
		return text
	}
	end, err := colorful.Hex(to)
	if err != nil {
		return text
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}
	if len(runes) == 1 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(from)).Render(text)
	}

	result := ""
	for index, char := range runes {
		blend := start.BlendHcl(end, float64(index)/float64(len(runes)-1)).Clamped()
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color(blend.Hex())).
			Bold(true).
			Render(string(char))
	}
	return result
}
