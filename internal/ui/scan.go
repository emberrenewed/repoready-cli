package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScanTask struct {
	Label string
	Run   func() error
}

type scanStatus int

const (
	scanWaiting scanStatus = iota
	scanRunning
	scanDone
	scanFailed
)

type scanTaskFinishedMsg struct {
	index int
	err   error
}

type scanModel struct {
	title    string
	tasks    []ScanTask
	statuses []scanStatus
	current  int
	frame    int
	spinner  spinner.Model
	progress progress.Model
	err      error
}

func newScanModel(title string, tasks []ScanTask) scanModel {
	statuses := make([]scanStatus, len(tasks))
	if len(statuses) > 0 {
		statuses[0] = scanRunning
	}
	return scanModel{
		title:    title,
		tasks:    tasks,
		statuses: statuses,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Points),
			spinner.WithStyle(PinkStyle),
		),
		progress: progress.New(
			progress.WithWidth(34),
			progress.WithGradient("#FF79C6", "#8BE9FD"),
		),
	}
}

func (m scanModel) Init() tea.Cmd {
	if len(m.tasks) == 0 {
		return tea.Quit
	}
	return tea.Batch(m.spinner.Tick, runScanTask(0, m.tasks[0]))
}

func (m scanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.frame++
		return m, cmd
	case scanTaskFinishedMsg:
		if msg.err != nil {
			m.statuses[msg.index] = scanFailed
			m.err = msg.err
			return m, tea.Quit
		}
		m.statuses[msg.index] = scanDone
		if msg.index == len(m.tasks)-1 {
			return m, tea.Quit
		}
		m.current = msg.index + 1
		m.statuses[m.current] = scanRunning
		return m, runScanTask(m.current, m.tasks[m.current])
	default:
		return m, nil
	}
}

func (m scanModel) View() string {
	var rows []string
	for index, task := range m.tasks {
		rows = append(rows, statusMark(m.statuses[index], m.spinner.View())+"  "+task.Label)
	}

	percent := 0.0
	if len(m.tasks) > 0 {
		completed := 0
		for _, status := range m.statuses {
			if status == scanDone {
				completed++
			}
		}
		percent = float64(completed) / float64(len(m.tasks))
	}

	sparkles := []string{"✦", "♡", "✿", "✨"}
	pulse := sparkles[m.frame%len(sparkles)]
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		GradientText(pulse+" "+m.title+" "+pulse, "#FF79C6", "#BD93F9"),
		MutedStyle.Render("reading the repository one clue at a time"),
		"",
		strings.Join(rows, "\n"),
		"",
		PinkStyle.Render("progress  ")+m.progress.ViewAs(percent),
	)
	return renderBox(content, Pink)
}

func statusMark(status scanStatus, liveSpinner string) string {
	switch status {
	case scanRunning:
		return liveSpinner
	case scanDone:
		return SuccessStyle.Render("♥")
	case scanFailed:
		return ErrorStyle.Render("×")
	default:
		return MutedStyle.Render("♡")
	}
}

func runScanTask(index int, task ScanTask) tea.Cmd {
	return func() tea.Msg {
		return scanTaskFinishedMsg{index: index, err: task.Run()}
	}
}

func RunScanSequence(title string, tasks []ScanTask) error {
	result, err := tea.NewProgram(newScanModel(title, tasks)).Run()
	if err != nil {
		return err
	}
	model, ok := result.(scanModel)
	if !ok {
		return fmt.Errorf("scan animation ended unexpectedly")
	}
	return model.err
}
