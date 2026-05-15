package system

import "runtime"

var installInstructions = map[string]map[string]string{
	"Git": {
		"linux":   "sudo apt install git",
		"darwin":  "brew install git",
		"windows": "Install Git from the official Git website.",
	},
	"Node.js": {
		"linux":   "Install Node.js using NodeSource or nvm.",
		"darwin":  "brew install node",
		"windows": "Install Node.js LTS from the official website.",
	},
	"npm": {
		"linux":   "npm ships with Node.js. Install Node.js first.",
		"darwin":  "npm ships with Node.js. Install Node.js first.",
		"windows": "npm ships with Node.js. Install Node.js first.",
	},
	"pnpm": {
		"linux":   "corepack enable pnpm",
		"darwin":  "corepack enable pnpm",
		"windows": "Run `corepack enable pnpm` after installing Node.js.",
	},
	"yarn": {
		"linux":   "corepack enable yarn",
		"darwin":  "corepack enable yarn",
		"windows": "Run `corepack enable yarn` after installing Node.js.",
	},
	"PHP": {
		"linux":   "sudo apt install php",
		"darwin":  "brew install php",
		"windows": "Install PHP from the official PHP for Windows package.",
	},
	"Composer": {
		"linux":   "sudo apt install composer",
		"darwin":  "brew install composer",
		"windows": "Install Composer Setup from the official website.",
	},
	"Python 3": {
		"linux":   "sudo apt install python3",
		"darwin":  "brew install python",
		"windows": "Install Python 3 from the official website.",
	},
	"Python": {
		"linux":   "Install Python 3 with `sudo apt install python3`.",
		"darwin":  "brew install python",
		"windows": "Install Python 3 from the official website.",
	},
	"pip": {
		"linux":   "sudo apt install python3-pip",
		"darwin":  "python3 -m ensurepip --upgrade",
		"windows": "pip ships with the official Python installer.",
	},
	"Go": {
		"linux":   "Install Go from the official Go downloads page.",
		"darwin":  "brew install go",
		"windows": "Install Go from the official Go downloads page.",
	},
	"Cargo": {
		"linux":   "Install Rust using rustup.",
		"darwin":  "brew install rustup-init",
		"windows": "Install Rust using rustup-init.exe.",
	},
	"Rust": {
		"linux":   "Install Rust using rustup.",
		"darwin":  "brew install rustup-init",
		"windows": "Install Rust using rustup-init.exe.",
	},
	"Docker": {
		"linux":   "Install Docker Engine from Docker's official instructions.",
		"darwin":  "Install Docker Desktop for macOS.",
		"windows": "Install Docker Desktop for Windows.",
	},
	"Docker Compose": {
		"linux":   "Install the Docker Compose plugin.",
		"darwin":  "Docker Desktop includes Docker Compose.",
		"windows": "Docker Desktop includes Docker Compose.",
	},
	"Java": {
		"linux":   "sudo apt install default-jdk",
		"darwin":  "brew install openjdk",
		"windows": "Install a JDK such as Temurin.",
	},
	"Maven": {
		"linux":   "sudo apt install maven",
		"darwin":  "brew install maven",
		"windows": "Install Maven from the official Apache website.",
	},
	"Gradle": {
		"linux":   "sudo apt install gradle",
		"darwin":  "brew install gradle",
		"windows": "Install Gradle from the official website.",
	},
	"Flutter": {
		"linux":   "Install Flutter from the official Flutter SDK guide.",
		"darwin":  "brew install --cask flutter",
		"windows": "Install Flutter from the official Flutter SDK guide.",
	},
	"Ruby": {
		"linux":   "sudo apt install ruby-full",
		"darwin":  "brew install ruby",
		"windows": "Install RubyInstaller for Windows.",
	},
	"Bundler": {
		"linux":   "gem install bundler",
		"darwin":  "gem install bundler",
		"windows": "gem install bundler",
	},
	".NET": {
		"linux":   "Install the .NET SDK from Microsoft's official guide.",
		"darwin":  "brew install --cask dotnet-sdk",
		"windows": "Install the .NET SDK from Microsoft's official website.",
	},
	"MySQL": {
		"linux":   "sudo apt install mysql-client",
		"darwin":  "brew install mysql-client",
		"windows": "Install MySQL from the official website.",
	},
	"PostgreSQL": {
		"linux":   "sudo apt install postgresql-client",
		"darwin":  "brew install libpq",
		"windows": "Install PostgreSQL from the official website.",
	},
	"MongoDB": {
		"linux":   "Install MongoDB Shell from MongoDB's official guide.",
		"darwin":  "brew install mongosh",
		"windows": "Install MongoDB Shell from MongoDB's official guide.",
	},
	"Redis": {
		"linux":   "sudo apt install redis-server",
		"darwin":  "brew install redis",
		"windows": "Use Memurai or WSL for Redis on Windows.",
	},
}

func InstructionFor(name string) string {
	byOS, ok := installInstructions[name]
	if !ok {
		return "See the official installation guide for this tool."
	}
	if instruction, ok := byOS[runtime.GOOS]; ok {
		return instruction
	}
	return "See the official installation guide for this tool."
}
