package prompt

import (
	"bytes"
	"strings"

	"github.com/nao1215/markdown"
)

func CodeAnalysis(code string) string {
	builder := new(bytes.Buffer)
	md := markdown.NewMarkdown(builder).
		H1("Code Analysis").
		H2("System").
		PlainText("You are a code analyzer. Your work is generating a explanation of a code block.").
		H2("Code Block").
		CodeBlocks("go", code).
		H2("Requirement").
		PlainText("Write a explanation of the code block above. The explanation must be written in '## explanation' block. You must not make other contents except for the explanation.")
	_ = md.Build()

	return UnwrapCodeAnalysis(md.String())
}

func UnwrapCodeAnalysis(response string) string {
	s := strings.Index(response, "## explanation")
	if s == -1 {
		s = strings.Index(response, "## Explanation")
		if s == -1 {
			return response
		}
	}

	return strings.TrimSpace(response[s+len("## explanation"):])
}
