package system

type ToolDefinition struct {
	Name    string
	Command string
	Args    []string
}

var toolDefinitions = []ToolDefinition{
	{Name: "Git", Command: "git", Args: []string{"--version"}},
	{Name: "Node.js", Command: "node", Args: []string{"--version"}},
	{Name: "npm", Command: "npm", Args: []string{"--version"}},
	{Name: "pnpm", Command: "pnpm", Args: []string{"--version"}},
	{Name: "yarn", Command: "yarn", Args: []string{"--version"}},
	{Name: "PHP", Command: "php", Args: []string{"--version"}},
	{Name: "Composer", Command: "composer", Args: []string{"--version"}},
	{Name: "Python", Command: "python", Args: []string{"--version"}},
	{Name: "Python 3", Command: "python3", Args: []string{"--version"}},
	{Name: "pip", Command: "pip", Args: []string{"--version"}},
	{Name: "pip3", Command: "pip3", Args: []string{"--version"}},
	{Name: "Go", Command: "go", Args: []string{"version"}},
	{Name: "Cargo", Command: "cargo", Args: []string{"--version"}},
	{Name: "Rust", Command: "rustc", Args: []string{"--version"}},
	{Name: "Docker", Command: "docker", Args: []string{"--version"}},
	{Name: "Docker Compose", Command: "docker", Args: []string{"compose", "version"}},
	{Name: "Java", Command: "java", Args: []string{"--version"}},
	{Name: "Maven", Command: "mvn", Args: []string{"--version"}},
	{Name: "Gradle", Command: "gradle", Args: []string{"--version"}},
	{Name: "Flutter", Command: "flutter", Args: []string{"--version"}},
	{Name: "Ruby", Command: "ruby", Args: []string{"--version"}},
	{Name: "Bundler", Command: "bundle", Args: []string{"--version"}},
	{Name: ".NET", Command: "dotnet", Args: []string{"--version"}},
	{Name: "MySQL", Command: "mysql", Args: []string{"--version"}},
	{Name: "PostgreSQL", Command: "psql", Args: []string{"--version"}},
	{Name: "MongoDB", Command: "mongosh", Args: []string{"--version"}},
	{Name: "Redis", Command: "redis-server", Args: []string{"--version"}},
}

func definitionByName(name string) (ToolDefinition, bool) {
	for _, definition := range toolDefinitions {
		if definition.Name == name {
			return definition, true
		}
	}
	return ToolDefinition{}, false
}
