package prompt

import (
	"bytes"
	"fmt"

	"github.com/nao1215/markdown"
)

func CodeAnalysis(code string) (string, error) {
	builder := new(bytes.Buffer)
	md := markdown.NewMarkdown(builder).
		H1("Code Analysis").
		H2("Description").
		PlainText("I want you to work as a senior software engineer at our company. Your work is code analysis our code and make summary.").
		H2("Code Block").
		CodeBlocks("go", code).
		H2("Requirement").
		PlainText("Please write a summary of the code block above. The first 100 characters are required. The summary should be written in <summary> tag.")
	if err := md.Build(); err != nil {
		return "", fmt.Errorf("failed to build markdown: %w", err)
	}

	return builder.String(), nil
}
