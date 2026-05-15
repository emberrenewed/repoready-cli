package system

import (
	"runtime"
	"strings"

	"repoready/internal/models"
)

type FixStep struct {
	Tool          models.ToolStatus
	Command       string
	CanRun        bool
	RequiresAdmin bool
}

var installCommands = map[string]map[string]string{
	"Git": {
		"linux":  "sudo apt install git",
		"darwin": "brew install git",
	},
	"pnpm": {
		"linux":  "corepack enable pnpm",
		"darwin": "corepack enable pnpm",
	},
	"yarn": {
		"linux":  "corepack enable yarn",
		"darwin": "corepack enable yarn",
	},
	"PHP": {
		"linux":  "sudo apt install php",
		"darwin": "brew install php",
	},
	"Composer": {
		"linux":  "sudo apt install composer",
		"darwin": "brew install composer",
	},
	"Python 3": {
		"linux":  "sudo apt install python3",
		"darwin": "brew install python",
	},
	"pip": {
		"linux":  "sudo apt install python3-pip",
		"darwin": "python3 -m ensurepip --upgrade",
	},
	"Maven": {
		"linux":  "sudo apt install maven",
		"darwin": "brew install maven",
	},
	"Gradle": {
		"linux":  "sudo apt install gradle",
		"darwin": "brew install gradle",
	},
	"Ruby": {
		"linux":  "sudo apt install ruby-full",
		"darwin": "brew install ruby",
	},
	"Bundler": {
		"linux":  "gem install bundler",
		"darwin": "gem install bundler",
	},
	"MySQL": {
		"linux":  "sudo apt install mysql-client",
		"darwin": "brew install mysql-client",
	},
	"PostgreSQL": {
		"linux":  "sudo apt install postgresql-client",
		"darwin": "brew install libpq",
	},
	"MongoDB": {
		"darwin": "brew install mongosh",
	},
	"Redis": {
		"linux":  "sudo apt install redis-server",
		"darwin": "brew install redis",
	},
}

func BuildFixPlan(tools []models.ToolStatus) []FixStep {
	steps := make([]FixStep, 0, len(tools))
	for _, tool := range tools {
		command, canRun := installCommand(tool.Name)
		steps = append(steps, FixStep{
			Tool:          tool,
			Command:       command,
			CanRun:        canRun,
			RequiresAdmin: strings.HasPrefix(command, "sudo "),
		})
	}
	return steps
}

func installCommand(name string) (string, bool) {
	byOS, ok := installCommands[name]
	if !ok {
		return "", false
	}
	command, ok := byOS[runtime.GOOS]
	return command, ok
}
