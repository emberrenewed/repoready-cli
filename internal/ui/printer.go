package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"repoready/internal/models"
	"repoready/internal/system"
)

func PrintBanner() {
	PrintAnimatedBanner()
}

func PrintFooter() {
	fmt.Println()
	fmt.Println(Footer())
}

func Section(title string) {
	fmt.Println(SectionStyle.Render("✦ " + title))
	fmt.Println()
}

func PrintProjectAnalysis(project models.ProjectInfo) {
	Section("Project Scan")
	fmt.Println(ProjectTable(project))
	fmt.Println()
}

func PrintSystemCheck(report models.SystemReport) {
	Section("Your System")
	fmt.Println(InfoBox(fmt.Sprintf("Host system: %s/%s\nChecking only the tools this repository actually needs.", report.OS, report.Architecture)))
	fmt.Println()
	fmt.Println(ToolTable(report))
	fmt.Println()
}

func PrintDiagnosis(project models.ProjectInfo, report models.SystemReport) []models.ToolStatus {
	missing := system.MissingTools(report)
	Section("Diagnosis")
	if project.MainLanguage == "Unknown" {
		RunMomentAnimation("thinking")
		fmt.Println(WarningBox("🌙 UNKNOWN STACK\n\nI scanned the repository files, but I could not confidently detect a supported language yet.\n\nCreated by CodeAhmad"))
		return missing
	}
	if len(missing) == 0 {
		RunMomentAnimation("ready")
		fmt.Println(SuccessBox(fmt.Sprintf(
			"🌸 READY\n\nRepository language: %s\nEvidence: %s\nYour system already has everything this repository needs.\n\nCreated by CodeAhmad",
			project.MainLanguage,
			evidenceText(project.LanguageEvidence),
		)))
		return nil
	}

	RunMomentAnimation("warning")
	runtimeName := system.PrimaryRuntimeForLanguage(project.MainLanguage)
	names := toolNames(missing)
	var message string
	if runtimeName != "" && hasTool(missing, runtimeName) {
		message = fmt.Sprintf(
			"⚠️ PROBLEM FOUND\n\nRepository language: %s\nEvidence: %s\nYour system does not have %s installed.\nYou must download/install: %s\n\nCreated by CodeAhmad",
			project.MainLanguage,
			evidenceText(project.LanguageEvidence),
			runtimeName,
			strings.Join(names, ", "),
		)
	} else {
		message = fmt.Sprintf(
			"⚠️ PROBLEM FOUND\n\nRepository language: %s\nEvidence: %s\nYour system is missing: %s\nYou must download/install: %s\n\nCreated by CodeAhmad",
			project.MainLanguage,
			evidenceText(project.LanguageEvidence),
			strings.Join(names, ", "),
			strings.Join(names, ", "),
		)
	}
	fmt.Println(WarningBox(message))
	fmt.Println()
	fmt.Println(MissingToolsTable(missing))
	return missing
}

func PrintError(message, fix string) {
	content := "💔 ERROR\n\n" + message
	if strings.TrimSpace(fix) != "" {
		content += "\n\nFix: " + fix
	}
	content += "\n\nCreated by CodeAhmad"
	fmt.Println(ErrorBox(content))
}

func InfoBox(content string) string {
	return renderBox(content, Blue)
}

func SuccessBox(content string) string {
	return renderBox(content, Green)
}

func WarningBox(content string) string {
	return renderBox(content, Yellow)
}

func ErrorBox(content string) string {
	return renderBox(content, Red)
}

func HelpScreen() string {
	body := lipgloss.JoinVertical(
		lipgloss.Center,
		GradientText("🌸 RepoReady", "#FF79C6", "#BD93F9"),
		SubtitleStyle.Render("Cute GitHub project assistant"),
		MutedStyle.Render("Created by CodeAhmad"),
		"",
		PinkStyle.Render("Usage"),
		"  repoready https://github.com/user/repo",
		"  repoready",
		"",
		BlueStyle.Render("Scan a GitHub repo, compare your system, and help with the next step."),
	)
	return BoxStyle.Copy().Width(66).Align(lipgloss.Center).Render(body) + "\n\n" + Footer()
}

func renderBox(content string, border lipgloss.TerminalColor) string {
	wrapped := wordwrap.String(content, 62)
	return BoxStyle.Copy().Width(68).BorderForeground(border).Render(wrapped)
}

func PrintFixPlan(steps []system.FixStep) {
	Section("Gentle Fix Plan")
	fmt.Println(InfoBox("I can prepare the fix now.\nCommands that need administrator access are shown, not run automatically."))
	fmt.Println()
	fmt.Println(FixPlanTable(steps))
}

func PrintFixResult(report models.SystemReport) {
	Section("Fix Result")
	missing := system.MissingTools(report)
	if len(missing) == 0 {
		RunMomentAnimation("ready")
		fmt.Println(SuccessBox("🌸 FIXED\n\nThe required tools are now available.\n\nCreated by CodeAhmad"))
		return
	}
	RunMomentAnimation("warning")
	fmt.Println(WarningBox("⚠️ STILL MISSING\n\nSome tools still need manual installation before this repository can run.\n\nCreated by CodeAhmad"))
	fmt.Println()
	fmt.Println(MissingToolsTable(missing))
}

func PrintDownloadOffer(path string) {
	Section("Ready To Download")
	fmt.Println(InfoBox("Your system is ready for this project.\n\nIf you continue, RepoReady will download it to:\n" + path))
}

func PrintDownloadSuccess(path string) {
	Section("Download Complete")
	RunMomentAnimation("download")
	fmt.Println(SuccessBox("📦 DOWNLOADED\n\nProject saved to:\n" + path + "\n\nCreated by CodeAhmad"))
}

func toolNames(tools []models.ToolStatus) []string {
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	return names
}

func hasTool(tools []models.ToolStatus, name string) bool {
	for _, tool := range tools {
		if tool.Name == name {
			return true
		}
	}
	return false
}

func evidenceText(values []string) string {
	if len(values) == 0 {
		return "source files"
	}
	return strings.Join(values, ", ")
}
