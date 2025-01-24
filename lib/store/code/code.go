package code

import "context"

type Data struct {
	Id          int
	CodeBlock   string
	FilePath    string
	Line        int
	Description string
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
