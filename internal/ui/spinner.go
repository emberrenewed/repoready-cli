package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type taskFinishedMsg struct {
	err error
}

type spinnerModel struct {
	spinner spinner.Model
	label   string
	done    <-chan error
	err     error
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForTask(m.done))
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case taskFinishedMsg:
		m.err = msg.err
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m spinnerModel) View() string {
	return m.spinner.View() + " " + SubtitleStyle.Render(m.label)
}

func waitForTask(done <-chan error) tea.Cmd {
	return func() tea.Msg {
		return taskFinishedMsg{err: <-done}
	}
}

func RunSpinner(label string, task func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- task()
	}()

	model := spinnerModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Points),
			spinner.WithStyle(PinkStyle),
		),
		label: label,
		done:  done,
	}
	result, err := tea.NewProgram(model).Run()
	if err != nil {
		return err
	}
	return result.(spinnerModel).err
}

func ProgressLine(label string, percent float64) string {
	bar := progress.New(
		progress.WithWidth(28),
		progress.WithGradient("#FF79C6", "#8BE9FD"),
	)
	return PinkStyle.Render(label+" ") + bar.ViewAs(percent)
}
