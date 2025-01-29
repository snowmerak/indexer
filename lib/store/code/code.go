package code

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/nao1215/markdown"
	"github.com/snowmerak/indexer/pkg/config"
)

type Data struct {
	Id          int
	CodeBlock   string
	FilePath    string
	Line        int
	Description string
}

func (d *Data) ToMarkdown(language string) string {
	builder := strings.Builder{}
	md := markdown.NewMarkdown(&builder).
		H2(fmt.Sprintf("Snippet %d", d.Id)).
		H3("Code Block").
		CodeBlocks(markdown.SyntaxHighlight(language), d.CodeBlock).
		H3("File Path").
		BulletList(fmt.Sprintf("File Path: %s", d.FilePath), fmt.Sprintf("Line: %d", d.Line)).
		H3("Description").
		PlainText(d.Description)
	_ = md.Build()

	return builder.String()
}

type Store interface {
	Create(ctx context.Context) error
	Save(ctx context.Context, id int, codeBlock string, filePath string, line int, description string) error
	Get(ctx context.Context, id int) (*Data, error)
	Gets(ctx context.Context, ids ...int) ([]*Data, error)
	Delete(ctx context.Context, id int) error
	Deletes(ctx context.Context, ids ...int) error
	Drop(ctx context.Context) error
}

var registeredStore = sync.Map{}

type StoreConstructor func(*config.ClientConfig) (Store, error)

func RegisterStore(name string, store StoreConstructor) {
	registeredStore.Store(name, store)
}

func GetStore(name string, config *config.ClientConfig) (Store, error) {
	if v, ok := registeredStore.Load(name); ok {
		store, err := v.(StoreConstructor)(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create store: %w", err)
		}
		return store, nil
	}
	return nil, fmt.Errorf("store not found: %s", name)
}
