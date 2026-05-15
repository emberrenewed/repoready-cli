package system

func PrimaryRuntimeForLanguage(language string) string {
	switch language {
	case "JavaScript", "TypeScript":
		return "Node.js"
	case "PHP":
		return "PHP"
	case "Python":
		return "Python 3"
	case "Go":
		return "Go"
	case "Rust":
		return "Rust"
	case "Java":
		return "Java"
	case "Ruby":
		return "Ruby"
	case "Dart":
		return "Flutter"
	case "C#":
		return ".NET"
	default:
		return ""
	}
}
