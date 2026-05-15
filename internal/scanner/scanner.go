package scanner

import "repoready/internal/models"

func Analyze(root, repoURL, owner, repoName string) (models.ProjectInfo, error) {
	files, err := newFileSet(root)
	if err != nil {
		return models.ProjectInfo{}, err
	}

	packageJSON := readPackageJSON(files.path("package.json"))
	composerJSON := readComposerJSON(files.path("composer.json"))
	dependencyText := readDependencyText(files)

	info := models.ProjectInfo{
		RepoURL:          repoURL,
		RepoName:         repoName,
		Owner:            owner,
		LocalPath:        root,
		MainLanguage:     detectMainLanguage(files, packageJSON),
		LanguageEvidence: detectLanguageEvidence(files),
		Frameworks:       detectFrameworks(files, packageJSON, composerJSON, dependencyText),
		PackageManagers:  detectPackageManagers(files),
		DetectedFiles:    files.detectedFiles(),
		Databases:        detectDatabases(files),
		NodeScripts:      packageJSON.Scripts,
		HasEnvExample:    files.has(".env.example"),
		HasEnv:           files.has(".env"),
		HasDocker:        files.hasAny("Dockerfile", "docker-compose.yml", "compose.yml"),
		HasReadme:        files.has("README.md"),
	}
	info.RequiredTools = detectRequiredTools(&info)
	return info, nil
}
