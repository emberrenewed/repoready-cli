package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type packageJSONFile struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

func (p packageJSONFile) hasDependency(name string) bool {
	_, direct := p.Dependencies[name]
	_, dev := p.DevDependencies[name]
	return direct || dev
}

type composerJSONFile struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

func (c composerJSONFile) hasDependency(name string) bool {
	_, direct := c.Require[name]
	_, dev := c.RequireDev[name]
	return direct || dev
}

type pyProjectFile struct {
	Project struct {
		Dependencies []string `toml:"dependencies"`
	} `toml:"project"`
	Tool map[string]any `toml:"tool"`
}

type fileSet struct {
	root            string
	files           map[string]string
	directories     map[string]struct{}
	allFiles        []string
	extensionCounts map[string]int
}

func newFileSet(root string) (fileSet, error) {
	set := fileSet{
		root:            root,
		files:           make(map[string]string),
		directories:     make(map[string]struct{}),
		extensionCounts: make(map[string]int),
	}

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}

		if entry.IsDir() {
			name := entry.Name()
			if shouldSkipDirectory(name) {
				return filepath.SkipDir
			}
			set.directories[name] = struct{}{}
			return nil
		}

		base := entry.Name()
		if existing, ok := set.files[base]; !ok || isShallower(rel, existing) {
			set.files[base] = rel
		}
		set.allFiles = append(set.allFiles, rel)
		ext := strings.ToLower(filepath.Ext(base))
		if ext != "" {
			set.extensionCounts[ext]++
		}
		return nil
	})
	if err != nil {
		return fileSet{}, err
	}
	return set, nil
}

func (f fileSet) has(name string) bool {
	_, ok := f.files[name]
	return ok
}

func (f fileSet) hasAny(names ...string) bool {
	for _, name := range names {
		if f.has(name) {
			return true
		}
	}
	return false
}

func (f fileSet) hasDirectory(name string) bool {
	_, ok := f.directories[name]
	return ok
}

func (f fileSet) hasExtension(ext string) bool {
	return f.extensionCounts[strings.ToLower(ext)] > 0
}

func (f fileSet) countExtension(ext string) int {
	return f.extensionCounts[strings.ToLower(ext)]
}

func (f fileSet) path(name string) string {
	rel, ok := f.files[name]
	if !ok {
		return filepath.Join(f.root, name)
	}
	return filepath.Join(f.root, rel)
}

func (f fileSet) detectedFiles() []string {
	targets := []string{
		"package.json",
		"pnpm-lock.yaml",
		"yarn.lock",
		"package-lock.json",
		"composer.json",
		"artisan",
		"requirements.txt",
		"pyproject.toml",
		"Pipfile",
		"go.mod",
		"Cargo.toml",
		"pom.xml",
		"build.gradle",
		"gradlew",
		"Gemfile",
		"Dockerfile",
		"docker-compose.yml",
		"compose.yml",
		"pubspec.yaml",
		".env.example",
		".env",
		"README.md",
		"tsconfig.json",
	}

	var found []string
	for _, target := range targets {
		if rel, ok := f.files[target]; ok {
			found = append(found, rel)
		}
	}
	sort.Strings(found)
	return found
}

func readPackageJSON(path string) packageJSONFile {
	var parsed packageJSONFile
	data, err := os.ReadFile(path)
	if err != nil {
		return parsed
	}
	_ = json.Unmarshal(data, &parsed)
	return parsed
}

func readComposerJSON(path string) composerJSONFile {
	var parsed composerJSONFile
	data, err := os.ReadFile(path)
	if err != nil {
		return parsed
	}
	_ = json.Unmarshal(data, &parsed)
	return parsed
}

func readDependencyText(files fileSet) string {
	var chunks []string
	for _, name := range []string{"requirements.txt", "Pipfile", "Gemfile", "pom.xml", "build.gradle"} {
		if files.has(name) {
			if data, err := os.ReadFile(files.path(name)); err == nil {
				chunks = append(chunks, string(data))
			}
		}
	}

	if files.has("pyproject.toml") {
		var project pyProjectFile
		if data, err := os.ReadFile(files.path("pyproject.toml")); err == nil {
			if toml.Unmarshal(data, &project) == nil {
				chunks = append(chunks, strings.Join(project.Project.Dependencies, "\n"))
			}
			chunks = append(chunks, string(data))
		}
	}

	return strings.Join(chunks, "\n")
}

func detectPackageManagers(files fileSet) []string {
	var managers []string

	if files.has("composer.json") {
		managers = appendUnique(managers, "Composer")
	}
	if files.has("package.json") {
		switch {
		case files.has("pnpm-lock.yaml"):
			managers = appendUnique(managers, "pnpm")
		case files.has("yarn.lock"):
			managers = appendUnique(managers, "yarn")
		default:
			managers = appendUnique(managers, "npm")
		}
	}
	if files.has("requirements.txt") || files.has("pyproject.toml") {
		managers = appendUnique(managers, "pip")
	}
	if files.has("Pipfile") {
		managers = appendUnique(managers, "Pipenv")
	}
	if files.has("go.mod") {
		managers = appendUnique(managers, "Go Modules")
	}
	if files.has("Cargo.toml") {
		managers = appendUnique(managers, "Cargo")
	}
	if files.has("pom.xml") {
		managers = appendUnique(managers, "Maven")
	}
	if files.has("build.gradle") || files.has("gradlew") {
		managers = appendUnique(managers, "Gradle")
	}
	if files.has("Gemfile") {
		managers = appendUnique(managers, "Bundler")
	}
	if files.has("pubspec.yaml") {
		managers = appendUnique(managers, "Flutter Pub")
	}

	return managers
}

func detectDatabases(files fileSet) []string {
	var text strings.Builder
	for _, name := range []string{".env.example", "docker-compose.yml", "compose.yml"} {
		if files.has(name) {
			if data, err := os.ReadFile(files.path(name)); err == nil {
				text.Write(data)
				text.WriteByte('\n')
			}
		}
	}

	lower := strings.ToLower(text.String())
	var databases []string
	if strings.Contains(lower, "mysql") || strings.Contains(lower, "mariadb") {
		databases = appendUnique(databases, "MySQL")
	}
	if strings.Contains(lower, "postgres") || strings.Contains(lower, "pgsql") {
		databases = appendUnique(databases, "PostgreSQL")
	}
	if strings.Contains(lower, "mongo") {
		databases = appendUnique(databases, "MongoDB")
	}
	if strings.Contains(lower, "redis") {
		databases = appendUnique(databases, "Redis")
	}
	return databases
}

func appendUnique(values []string, items ...string) []string {
	existing := make(map[string]struct{}, len(values))
	for _, value := range values {
		existing[value] = struct{}{}
	}
	for _, item := range items {
		if _, ok := existing[item]; ok {
			continue
		}
		values = append(values, item)
		existing[item] = struct{}{}
	}
	return values
}

func shouldSkipDirectory(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", "dist", "build", ".next", "target", "venv", ".venv":
		return true
	default:
		return false
	}
}

func isShallower(candidate, current string) bool {
	return strings.Count(candidate, string(os.PathSeparator)) < strings.Count(current, string(os.PathSeparator))
}

func detectLanguageEvidence(files fileSet) []string {
	type evidenceRule struct {
		name string
		ok   bool
	}

	rules := []evidenceRule{
		{name: "pubspec.yaml", ok: files.has("pubspec.yaml")},
		{name: "go.mod", ok: files.has("go.mod")},
		{name: "Cargo.toml", ok: files.has("Cargo.toml")},
		{name: "requirements.txt", ok: files.has("requirements.txt")},
		{name: "pyproject.toml", ok: files.has("pyproject.toml")},
		{name: "Pipfile", ok: files.has("Pipfile")},
		{name: "composer.json", ok: files.has("composer.json")},
		{name: "Gemfile", ok: files.has("Gemfile")},
		{name: "pom.xml", ok: files.has("pom.xml")},
		{name: "build.gradle", ok: files.has("build.gradle")},
		{name: "package.json", ok: files.has("package.json")},
		{name: "tsconfig.json", ok: files.has("tsconfig.json")},
	}

	var evidence []string
	for _, rule := range rules {
		if rule.ok {
			evidence = append(evidence, rule.name)
		}
	}
	if len(evidence) > 0 {
		return evidence
	}

	sourceEvidence := map[string]string{
		".ts":   ".ts files",
		".tsx":  ".tsx files",
		".js":   ".js files",
		".jsx":  ".jsx files",
		".php":  ".php files",
		".py":   ".py files",
		".go":   ".go files",
		".rs":   ".rs files",
		".java": ".java files",
		".rb":   ".rb files",
		".dart": ".dart files",
		".cs":   ".cs files",
	}
	for extension, label := range sourceEvidence {
		if count := files.countExtension(extension); count > 0 {
			evidence = append(evidence, label+" ("+strconv.Itoa(count)+")")
		}
	}
	sort.Strings(evidence)
	return evidence
}
