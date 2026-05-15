package models

type ProjectInfo struct {
	RepoURL          string
	RepoName         string
	Owner            string
	LocalPath        string
	MainLanguage     string
	LanguageEvidence []string
	Frameworks       []string
	PackageManagers  []string
	RequiredTools    []string
	DetectedFiles    []string
	Databases        []string
	NodeScripts      map[string]string
	HasEnvExample    bool
	HasEnv           bool
	HasDocker        bool
	HasReadme        bool
}
