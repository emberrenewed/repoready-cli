package system

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"repoready/internal/models"
)

type Checker struct{}

func NewChecker() Checker {
	return Checker{}
}

func (Checker) CheckAll(ctx context.Context) models.SystemReport {
	return check(ctx, nil)
}

func (Checker) CheckRequired(ctx context.Context, required []string) models.SystemReport {
	requiredSet := make(map[string]struct{}, len(required))
	for _, name := range required {
		requiredSet[name] = struct{}{}
	}
	return check(ctx, requiredSet)
}

func check(ctx context.Context, requiredSet map[string]struct{}) models.SystemReport {
	report := models.SystemReport{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	for _, definition := range toolDefinitions {
		_, required := requiredSet[definition.Name]
		if requiredSet != nil && !required {
			continue
		}
		report.Tools = append(report.Tools, inspectTool(ctx, definition, required))
	}
	return report
}

func inspectTool(ctx context.Context, definition ToolDefinition, required bool) models.ToolStatus {
	status := models.ToolStatus{
		Name:                definition.Name,
		Command:             strings.TrimSpace(strings.Join(append([]string{definition.Command}, definition.Args...), " ")),
		Required:            required,
		InstallInstructions: InstructionFor(definition.Name),
		Version:             "-",
	}

	cmd := exec.CommandContext(ctx, definition.Command, definition.Args...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		status.Installed = true
		status.Version = firstLine(string(output))
		return status
	}

	for _, alternative := range alternativesFor(definition.Name) {
		cmd := exec.CommandContext(ctx, alternative.Command, alternative.Args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}
		version := firstLine(string(output))
		if alternative.MustContain != "" && !strings.Contains(strings.ToLower(version), strings.ToLower(alternative.MustContain)) {
			continue
		}
		status.Installed = true
		status.Command = strings.TrimSpace(strings.Join(append([]string{alternative.Command}, alternative.Args...), " "))
		status.Version = version
		return status
	}

	return status
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	lines := strings.Split(value, "\n")
	return strings.TrimSpace(lines[0])
}

func MissingTools(report models.SystemReport) []models.ToolStatus {
	var missing []models.ToolStatus
	for _, tool := range report.Tools {
		if !tool.Installed {
			missing = append(missing, tool)
		}
	}
	return missing
}

func ToolInstalled(report models.SystemReport, name string) bool {
	for _, tool := range report.Tools {
		if tool.Name == name {
			return tool.Installed
		}
	}
	return false
}

func Definition(name string) (ToolDefinition, bool) {
	return definitionByName(name)
}

type toolAlternative struct {
	Command     string
	Args        []string
	MustContain string
}

func alternativesFor(name string) []toolAlternative {
	switch name {
	case "pip":
		return []toolAlternative{
			{Command: "pip3", Args: []string{"--version"}},
			{Command: "python3", Args: []string{"-m", "pip", "--version"}},
		}
	case "Python 3":
		return []toolAlternative{
			{Command: "python", Args: []string{"--version"}, MustContain: "python 3"},
		}
	default:
		return nil
	}
}
