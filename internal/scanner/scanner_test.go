package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeLaravelProject(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "composer.json"), `{"require":{"laravel/framework":"^12.0"}}`)
	writeFile(t, filepath.Join(root, "artisan"), "")
	writeFile(t, filepath.Join(root, "package.json"), `{"scripts":{"dev":"vite"},"devDependencies":{"vite":"latest"}}`)
	writeFile(t, filepath.Join(root, "package-lock.json"), "")
	writeFile(t, filepath.Join(root, ".env.example"), "DB_CONNECTION=mysql")

	project, err := Analyze(root, "https://github.com/codeahmad/demo", "codeahmad", "demo")
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if project.MainLanguage != "PHP" {
		t.Fatalf("expected PHP, got %s", project.MainLanguage)
	}
	assertContains(t, project.Frameworks, "Laravel")
	assertContains(t, project.Frameworks, "Vite")
	assertContains(t, project.PackageManagers, "Composer")
	assertContains(t, project.PackageManagers, "npm")
	assertContains(t, project.Databases, "MySQL")
}

func TestAnalyzeNodeProjectDetectsTypeScriptAndPnpm(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "package.json"), `{"dependencies":{"react":"latest"},"devDependencies":{"typescript":"latest"},"scripts":{"dev":"vite"}}`)
	writeFile(t, filepath.Join(root, "pnpm-lock.yaml"), "")

	project, err := Analyze(root, "https://github.com/codeahmad/web", "codeahmad", "web")
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if project.MainLanguage != "TypeScript" {
		t.Fatalf("expected TypeScript, got %s", project.MainLanguage)
	}
	assertContains(t, project.Frameworks, "React")
	assertContains(t, project.PackageManagers, "pnpm")
}

func TestAnalyzeDetectsLanguageFromNestedSourceFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src", "service"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(root, "src", "service", "main.py"), "print('hello')")

	project, err := Analyze(root, "https://github.com/codeahmad/python", "codeahmad", "python")
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if project.MainLanguage != "Python" {
		t.Fatalf("expected Python, got %s", project.MainLanguage)
	}
	assertContains(t, project.RequiredTools, "Python 3")
	assertContains(t, project.LanguageEvidence, ".py files (1)")
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func assertContains(t *testing.T, values []string, target string) {
	t.Helper()
	for _, value := range values {
		if value == target {
			return
		}
	}
	t.Fatalf("expected %q in %v", target, values)
}
