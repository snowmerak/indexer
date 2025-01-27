package prompt

import (
	"bytes"
	"strings"

	"github.com/nao1215/markdown"
)

func CodeAnalysis(language string, code string) string {
	builder := new(bytes.Buffer)
	md := markdown.NewMarkdown(builder).
		H1("Code Analysis").
		H2("System Prompt").
		PlainText("I want you to act as senior software engineer who make review and documentation.").
		H2("Code Block").
		CodeBlocks(markdown.SyntaxHighlight(language), code).
		H2("Requirement").
		PlainText(`Analyze the following code and provide a detailed, step-by-step explanation of its behavior.
Specifically, identify each variable used in the code, describe its role, data type (if discernible), and track how its value changes during execution.
Focus on explaining how the code actually works, rather than just its structure.`)
	_ = md.Build()

	return UnwrapCodeAnalysis(md.String())
}

func UnwrapCodeAnalysis(response string) string {
	s := strings.Index(response, "## description")
	if s == -1 {
		s = strings.Index(response, "## Description")
		if s == -1 {
			return response
		}
	}

	return strings.TrimSpace(response[s+len("## description"):])
}
