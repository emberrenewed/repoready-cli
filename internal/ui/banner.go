package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Banner() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		GradientText("🌸 RepoReady", "#FF79C6", "#BD93F9"),
		SubtitleStyle.Render("Cute GitHub project assistant"),
		MutedStyle.Render("Created by emberrenewed"),
	)
	return BoxStyle.Copy().
		Width(48).
		Align(lipgloss.Center).
		Render(content)
}

func Footer() string {
	line := strings.Repeat("─", 36)
	return FooterStyle.Render(line + "\nMade with 💜 by emberrenewed")
}

type introTickMsg time.Time

type introModel struct {
	frame int
}

func (m introModel) Init() tea.Cmd {
	return tea.Tick(95*time.Millisecond, func(t time.Time) tea.Msg {
		return introTickMsg(t)
	})
}

func (m introModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case introTickMsg:
		m.frame++
		if m.frame >= 6 {
			return m, tea.Quit
		}
		return m, tea.Tick(95*time.Millisecond, func(t time.Time) tea.Msg {
			return introTickMsg(t)
		})
	default:
		return m, nil
	}
}

func (m introModel) View() string {
	flowers := []string{"✦", "♡", "✿", "🌸", "✨", "🌸"}
	subtitles := []string{
		"waking up the scanner",
		"reading tiny clues",
		"matching your system",
		"preparing the magic",
		"Cute GitHub project assistant",
		"Cute GitHub project assistant",
	}
	frame := m.frame
	if frame >= len(flowers) {
		frame = len(flowers) - 1
	}
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		GradientText(flowers[frame]+" RepoReady", "#FF79C6", "#BD93F9"),
		SubtitleStyle.Render(subtitles[frame]),
		MutedStyle.Render("Created by emberrenewed"),
	)
	return BoxStyle.Copy().
		Width(48).
		Align(lipgloss.Center).
		Render(content)
}

func PrintAnimatedBanner() {
	_, _ = tea.NewProgram(introModel{}).Run()
	fmt.Println(Banner())
	fmt.Println()
}
