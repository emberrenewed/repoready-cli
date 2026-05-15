package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"repoready/internal/models"
	"repoready/internal/system"
)

func ProjectTable(project models.ProjectInfo) string {
	rows := [][]string{
		{"Repository", project.Owner + "/" + project.RepoName},
		{"Detected Language", fallback(project.MainLanguage)},
		{"Evidence", summarizeFiles(project.LanguageEvidence)},
		{"Framework", joinOrDash(project.Frameworks)},
		{"Package Manager", joinOrDash(project.PackageManagers)},
		{"Required Tools", joinOrDash(project.RequiredTools)},
		{"Database", joinOrDash(project.Databases)},
		{"Docker", yesNo(project.HasDocker)},
		{"Scanned Files", summarizeFiles(project.DetectedFiles)},
	}
	return renderTable([]string{"Field", "Value"}, rows)
}

func ToolTable(report models.SystemReport) string {
	var rows [][]string
	for _, tool := range report.Tools {
		status := WarningStyle.Render("❌ Missing")
		if tool.Installed {
			status = SuccessStyle.Render("✅ Available")
		}
		rows = append(rows, []string{tool.Name, status, fallback(tool.Version)})
	}
	return renderTable([]string{"Tool", "Status", "Version"}, rows)
}

func MissingToolsTable(tools []models.ToolStatus) string {
	var rows [][]string
	for _, tool := range tools {
		rows = append(rows, []string{tool.Name, tool.InstallInstructions})
	}
	return renderTable([]string{"Missing Tool", "Install Instruction"}, rows)
}

func FixPlanTable(steps []system.FixStep) string {
	var rows [][]string
	for _, step := range steps {
		mode := "manual"
		action := step.Tool.InstallInstructions
		if step.CanRun && !step.RequiresAdmin {
			mode = "can run"
			action = step.Command
		} else if step.CanRun {
			mode = "needs admin"
			action = step.Command
		}
		rows = append(rows, []string{step.Tool.Name, mode, action})
	}
	return renderTable([]string{"Tool", "Mode", "Fix"}, rows)
}

func renderTable(headers []string, rows [][]string) string {
	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(Pink)).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return TableHeaderStyle.Padding(0, 1)
			}
			style := lipgloss.NewStyle().Foreground(Text).Padding(0, 1)
			if col == 0 {
				style = style.Foreground(Purple).Bold(true)
			}
			if row%2 == 1 && col > 0 {
				style = style.Foreground(Text)
			}
			return style
		}).
		Render()
}

func joinOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func yesNo(value bool) string {
	if value {
		return "Found"
	}
	return "-"
}

func fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func summarizeFiles(files []string) string {
	if len(files) == 0 {
		return "-"
	}
	if len(files) <= 5 {
		return strings.Join(files, ", ")
	}
	return strings.Join(files[:5], ", ") + " +" + strconv.Itoa(len(files)-5) + " more"
}
