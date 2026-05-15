package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type momentTickMsg time.Time

type momentModel struct {
	kind  string
	frame int
}

func (m momentModel) Init() tea.Cmd {
	return tea.Tick(90*time.Millisecond, func(t time.Time) tea.Msg {
		return momentTickMsg(t)
	})
}

func (m momentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case momentTickMsg:
		m.frame++
		if m.frame >= 5 {
			return m, tea.Quit
		}
		return m, tea.Tick(90*time.Millisecond, func(t time.Time) tea.Msg {
			return momentTickMsg(t)
		})
	default:
		return m, nil
	}
}

func (m momentModel) View() string {
	frames, label, style := momentFrames(m.kind)
	frame := m.frame
	if frame >= len(frames) {
		frame = len(frames) - 1
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		style.Render(frames[frame]),
		" ",
		MutedStyle.Render(label),
	)
}

func momentFrames(kind string) ([]string, string, lipgloss.Style) {
	switch kind {
	case "ready":
		return []string{"♡", "✦", "🌸", "✨", "✅"}, "everything lines up beautifully", SuccessStyle
	case "download":
		return []string{"📦", "💜", "🌸", "✨", "✅"}, "project tucked safely into place", PinkStyle
	case "warning":
		return []string{"✦", "⚠", "⚠", "♡", "⚠"}, "a little care is needed here", WarningStyle
	default:
		return []string{"♡", "✦", "✿", "…", "♡"}, "still reading the clues", PurpleStyle
	}
}

func RunMomentAnimation(kind string) {
	_, _ = tea.NewProgram(momentModel{kind: kind}).Run()
}
