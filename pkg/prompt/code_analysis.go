package prompt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nao1215/markdown"
)

func CodeAnalysis(code string) (string, error) {
	builder := new(bytes.Buffer)
	md := markdown.NewMarkdown(builder).
		H1("Code Analysis").
		H2("System").
		PlainText("You are a code analyzer. Your work is generating a description of a code block.").
		H2("Code Block").
		CodeBlocks("go", code).
		H2("Requirement").
		PlainText("Write a description of the code block above. The description must be written in '## description' block. You must not make other contents except for the description.")
	if err := md.Build(); err != nil {
		return "", fmt.Errorf("failed to build markdown: %w", err)
	}

	return builder.String(), nil
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
