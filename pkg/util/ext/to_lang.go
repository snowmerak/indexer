package ext

func ToLang(ext string) string {
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".java":
		return "java"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".cs":
		return "csharp"
	case ".cpp":
		return "cpp"
	case ".rs":
		return "rust"
	case ".swift":
		return "swift"
	case ".kt":
		return "kotlin"
	case ".clj":
		return "clojure"
	case ".scala":
		return "scala"
	case ".r":
		return "r"
	case ".sh":
		return "shell"
	case ".pl":
		return "perl"
	case ".lua":
		return "lua"
	case ".dart":
		return "dart"
	case ".tsql":
		return "tsql"
	case ".v":
		return "verilog"
	case ".vhdl":
		return "vhdl"
	case ".asm":
		return "assembly"
	case ".sql":
		return "sql"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".xml":
		return "xml"
	case ".json":
		return "json"
	case ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".ini":
		return "ini"
	case ".md":
		return "markdown"
	default:
		return ""
	}
}
