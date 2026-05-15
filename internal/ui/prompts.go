package ui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type urlModel struct {
	input textinput.Model
	done  bool
}

func newURLModel() urlModel {
	input := textinput.New()
	input.Placeholder = "https://github.com/user/repo"
	input.Prompt = "🌸 repo url > "
	input.PromptStyle = PromptStyle
	input.TextStyle = lipgloss.NewStyle().Foreground(Text)
	input.PlaceholderStyle = MutedStyle
	input.CharLimit = 300
	input.Width = 52
	input.Focus()
	return urlModel{input: input}
}

func (m urlModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m urlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m urlModel) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		GradientText("Paste a GitHub repository URL", "#FF79C6", "#BD93F9"),
		"",
		m.input.View(),
	)
	return "\n" + BoxStyle.Copy().Width(64).Render(content) + "\n"
}

func PromptRepositoryURL() (string, error) {
	program := tea.NewProgram(newURLModel())
	result, err := program.Run()
	if err != nil {
		return "", err
	}
	model, ok := result.(urlModel)
	if !ok || !model.done {
		return "", errors.New("repository URL prompt cancelled")
	}
	return strings.TrimSpace(model.input.Value()), nil
}

func Confirm(question string, defaultYes bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		suffix := "[y/N]"
		if defaultYes {
			suffix = "[Y/n]"
		}
		fmt.Printf("%s %s ", PromptStyle.Render("💬 "+question), MutedStyle.Render(suffix))

		value, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		value = strings.TrimSpace(strings.ToLower(value))
		if value == "" {
			return defaultYes, nil
		}
		switch value {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println(WarningStyle.Render("⚠ answer with y or n"))
		}
	}
}
