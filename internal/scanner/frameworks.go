package scanner

import (
	"strings"
)

func detectFrameworks(files fileSet, packageJSON packageJSONFile, composerJSON composerJSONFile, dependencyText string) []string {
	var frameworks []string

	if packageJSON.hasDependency("react") {
		frameworks = appendUnique(frameworks, "React")
	}
	if packageJSON.hasDependency("vue") {
		frameworks = appendUnique(frameworks, "Vue")
	}
	if packageJSON.hasDependency("@angular/core") {
		frameworks = appendUnique(frameworks, "Angular")
	}
	if packageJSON.hasDependency("next") {
		frameworks = appendUnique(frameworks, "Next.js")
	}
	if packageJSON.hasDependency("vite") {
		frameworks = appendUnique(frameworks, "Vite")
	}
	if packageJSON.hasDependency("svelte") {
		frameworks = appendUnique(frameworks, "Svelte")
	}
	if packageJSON.hasDependency("express") {
		frameworks = appendUnique(frameworks, "Express")
	}
	if packageJSON.hasDependency("@nestjs/core") {
		frameworks = appendUnique(frameworks, "NestJS")
	}
	if packageJSON.hasDependency("react-native") {
		frameworks = appendUnique(frameworks, "React Native")
	}

	if files.has("artisan") && composerJSON.hasDependency("laravel/framework") {
		frameworks = appendUnique(frameworks, "Laravel")
	}

	lowerDeps := strings.ToLower(dependencyText)
	if strings.Contains(lowerDeps, "django") {
		frameworks = appendUnique(frameworks, "Django")
	}
	if strings.Contains(lowerDeps, "flask") {
		frameworks = appendUnique(frameworks, "Flask")
	}
	if strings.Contains(lowerDeps, "fastapi") {
		frameworks = appendUnique(frameworks, "FastAPI")
	}
	if strings.Contains(lowerDeps, "rails") && files.has("Gemfile") {
		frameworks = appendUnique(frameworks, "Rails")
	}
	if strings.Contains(lowerDeps, "spring-boot") {
		frameworks = appendUnique(frameworks, "Spring Boot")
	}
	if files.has("pubspec.yaml") {
		frameworks = appendUnique(frameworks, "Flutter")
	}
	if files.has("gradlew") || files.has("build.gradle") {
		frameworks = appendUnique(frameworks, "Android")
	}
	if files.hasDirectory("ios") {
		frameworks = appendUnique(frameworks, "iOS")
	}

	return frameworks
}
