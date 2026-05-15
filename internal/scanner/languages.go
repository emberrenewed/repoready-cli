package scanner

import "repoready/internal/models"

func detectMainLanguage(files fileSet, packageJSON packageJSONFile) string {
	switch {
	case files.has("pubspec.yaml"):
		return "Dart"
	case files.has("go.mod"):
		return "Go"
	case files.has("Cargo.toml"):
		return "Rust"
	case files.hasAny("requirements.txt", "pyproject.toml", "Pipfile"):
		return "Python"
	case files.has("composer.json"):
		return "PHP"
	case files.has("Gemfile"):
		return "Ruby"
	case files.hasAny("pom.xml", "build.gradle", "gradlew"):
		return "Java"
	case files.hasExtension(".csproj") || files.hasExtension(".sln"):
		return "C#"
	case files.has("package.json"):
		if files.has("tsconfig.json") || packageJSON.hasDependency("typescript") {
			return "TypeScript"
		}
		return "JavaScript"
	default:
		return detectLanguageFromSourceFiles(files)
	}
}

func detectRequiredTools(info *models.ProjectInfo) []string {
	var tools []string
	tools = appendUnique(tools, "Git")

	switch info.MainLanguage {
	case "JavaScript", "TypeScript":
		tools = appendUnique(tools, "Node.js")
	case "PHP":
		tools = appendUnique(tools, "PHP")
	case "Python":
		tools = appendUnique(tools, "Python 3")
	case "Go":
		tools = appendUnique(tools, "Go")
	case "Rust":
		tools = appendUnique(tools, "Rust", "Cargo")
	case "Java":
		tools = appendUnique(tools, "Java")
	case "Ruby":
		tools = appendUnique(tools, "Ruby")
	case "Dart":
		tools = appendUnique(tools, "Flutter")
	case "C#":
		tools = appendUnique(tools, ".NET")
	}

	for _, manager := range info.PackageManagers {
		switch manager {
		case "npm":
			tools = appendUnique(tools, "Node.js", "npm")
		case "pnpm":
			tools = appendUnique(tools, "Node.js", "pnpm")
		case "yarn":
			tools = appendUnique(tools, "Node.js", "yarn")
		case "Composer":
			tools = appendUnique(tools, "PHP", "Composer")
		case "pip":
			tools = appendUnique(tools, "Python 3", "pip")
		case "Pipenv":
			tools = appendUnique(tools, "Python 3", "pip")
		case "Go Modules":
			tools = appendUnique(tools, "Go")
		case "Cargo":
			tools = appendUnique(tools, "Cargo", "Rust")
		case "Maven":
			tools = appendUnique(tools, "Java", "Maven")
		case "Gradle":
			tools = appendUnique(tools, "Java", "Gradle")
		case "Bundler":
			tools = appendUnique(tools, "Ruby", "Bundler")
		case "Flutter Pub":
			tools = appendUnique(tools, "Flutter")
		}
	}

	if info.HasDocker {
		tools = appendUnique(tools, "Docker", "Docker Compose")
	}
	for _, database := range info.Databases {
		switch database {
		case "MySQL":
			tools = appendUnique(tools, "MySQL")
		case "PostgreSQL":
			tools = appendUnique(tools, "PostgreSQL")
		case "MongoDB":
			tools = appendUnique(tools, "MongoDB")
		case "Redis":
			tools = appendUnique(tools, "Redis")
		}
	}

	return tools
}

func detectLanguageFromSourceFiles(files fileSet) string {
	type candidate struct {
		language string
		score    int
	}

	candidates := []candidate{
		{language: "TypeScript", score: files.countExtension(".ts") + files.countExtension(".tsx")},
		{language: "JavaScript", score: files.countExtension(".js") + files.countExtension(".jsx")},
		{language: "PHP", score: files.countExtension(".php")},
		{language: "Python", score: files.countExtension(".py")},
		{language: "Go", score: files.countExtension(".go")},
		{language: "Rust", score: files.countExtension(".rs")},
		{language: "Java", score: files.countExtension(".java")},
		{language: "Ruby", score: files.countExtension(".rb")},
		{language: "Dart", score: files.countExtension(".dart")},
		{language: "C#", score: files.countExtension(".cs")},
	}

	best := candidate{language: "Unknown"}
	for _, candidate := range candidates {
		if candidate.score > best.score {
			best = candidate
		}
	}
	return best.language
}
