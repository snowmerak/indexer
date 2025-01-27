package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snowmerak/indexer/lib/store/code"
	queries2 "github.com/snowmerak/indexer/pkg/client/postgres/queries"
)

var _ code.Store = (*Store)(nil)

type Config struct {
	ConnectionString string

	TableName string
}

func NewConfig(connString string, tableName string) *Config {
	return &Config{
		ConnectionString: connString,
		TableName:        tableName,
	}
}

type Store struct {
	pool   *pgxpool.Pool
	config *Config
	conn   *queries2.Queries
}

func New(ctx context.Context, cfg *Config) (*Store, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	context.AfterFunc(ctx, func() {
		pool.Close()
	})

	return &Store{
		pool:   pool,
		config: cfg,
		conn:   queries2.New(pool),
	}, nil
}

//go:embed queries/schema.sql
var schema string

func (s *Store) Create(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, schema); err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (s *Store) Save(ctx context.Context, id int, codeBlock string, filePath string, line int, description string) error {
	if _, err := s.conn.CreateData(ctx, queries2.CreateDataParams{
		Project:     s.config.TableName,
		ID:          int64(id),
		CodeBlock:   pgtype.Text{String: codeBlock, Valid: true},
		FilePath:    pgtype.Text{String: filePath, Valid: true},
		Line:        pgtype.Int4{Int32: int32(line), Valid: true},
		Description: pgtype.Text{String: description, Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

func (s *Store) Get(ctx context.Context, id int) (*code.Data, error) {
	data, err := s.conn.GetData(ctx, queries2.GetDataParams{
		Project: s.config.TableName,
		ID:      int64(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	return &code.Data{
		Id:          int(data.ID),
		CodeBlock:   data.CodeBlock.String,
		FilePath:    data.FilePath.String,
		Line:        int(data.Line.Int32),
		Description: data.Description.String,
	}, nil
}

func (s *Store) Gets(ctx context.Context, ids ...int) ([]*code.Data, error) {
	i32l := make([]int32, len(ids))
	for i, id := range ids {
		i32l[i] = int32(id)
	}

	data, err := s.conn.GetDataList(ctx, queries2.GetDataListParams{
		Project: s.config.TableName,
		Column2: i32l,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	res := make([]*code.Data, len(data))
	for i, d := range data {
		res[i] = &code.Data{
			Id:          int(d.ID),
			CodeBlock:   d.CodeBlock.String,
			FilePath:    d.FilePath.String,
			Line:        int(d.Line.Int32),
			Description: d.Description.String,
		}
	}

	return res, nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	if _, err := s.conn.DeleteData(ctx, queries2.DeleteDataParams{
		Project: s.config.TableName,
		ID:      int64(id),
	}); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (s *Store) Deletes(ctx context.Context, ids ...int) error {
	i32l := make([]int32, len(ids))
	for i, id := range ids {
		i32l[i] = int32(id)
	}

	if _, err := s.conn.DeleteDataList(ctx, queries2.DeleteDataListParams{
		Project: s.config.TableName,
		Column2: i32l,
	}); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (s *Store) Drop(ctx context.Context) error {
	if _, err := s.conn.DeleteProjectData(ctx, s.config.TableName); err != nil {
		return fmt.Errorf("failed to drop data: %w", err)
	}

	return nil
}
