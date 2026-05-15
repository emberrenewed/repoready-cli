package system

import (
	"testing"

	"repoready/internal/models"
)

func TestPrimaryRuntimeForLanguage(t *testing.T) {
	if got := PrimaryRuntimeForLanguage("PHP"); got != "PHP" {
		t.Fatalf("expected PHP runtime, got %q", got)
	}
	if got := PrimaryRuntimeForLanguage("TypeScript"); got != "Node.js" {
		t.Fatalf("expected Node.js runtime, got %q", got)
	}
}

func TestBuildFixPlanKeepsMissingTools(t *testing.T) {
	tools := []models.ToolStatus{
		{Name: "Composer", InstallInstructions: "sudo apt install composer"},
		{Name: "Go", InstallInstructions: "Install Go manually."},
	}
	steps := BuildFixPlan(tools)
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if steps[0].Tool.Name != "Composer" || steps[1].Tool.Name != "Go" {
		t.Fatalf("unexpected steps: %#v", steps)
	}
}
