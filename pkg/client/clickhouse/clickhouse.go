package clickhouse

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/snowmerak/indexer/lib/store/code"
	"github.com/snowmerak/indexer/pkg/config"
)

type Config struct {
	tableName string

	addresses []string
	database  string
	username  string
	password  string

	dialTimeout     time.Duration
	maxOpenConn     int
	maxIdleConn     int
	connMaxLifetime time.Duration
	blockBufferSize int

	tlsConfig *tls.Config
}

func NewConfig(tableName string, addr ...string) *Config {
	return &Config{
		tableName: tableName,

		addresses: addr,
		username:  "",
		password:  "",
		database:  "",

		dialTimeout:     5 * time.Second,
		maxOpenConn:     10,
		maxIdleConn:     5,
		connMaxLifetime: 30 * time.Minute,
		blockBufferSize: 10,
	}
}

func (c *Config) WithDatabase(db string) *Config {
	c.database = db
	return c
}

func (c *Config) WithUsername(username string) *Config {
	c.username = username
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.password = password
	return c
}

func (c *Config) WithDialTimeout(timeout time.Duration) *Config {
	c.dialTimeout = timeout
	return c
}

func (c *Config) WithMaxOpenConn(max int) *Config {
	c.maxOpenConn = max
	return c
}

func (c *Config) WithMaxIdleConn(max int) *Config {
	c.maxIdleConn = max
	return c
}

func (c *Config) WithConnMaxLifetime(lifetime time.Duration) *Config {
	c.connMaxLifetime = lifetime
	return c
}

func (c *Config) WithBlockBufferSize(size int) *Config {
	c.blockBufferSize = size
	return c
}

func (c *Config) WithTLSConfig(tlsConfig *tls.Config) *Config {
	c.tlsConfig = tlsConfig
	return c
}

type Clickhouse struct {
	conn   driver.Conn
	config *Config
}

var _ code.Store = (*Clickhouse)(nil)

func init() {
	code.RegisterStore("clickhouse", func(ctx context.Context, cc *config.ClientConfig) (code.Store, error) {
		cfg := NewConfig(cc.Project, cc.Host...).
			WithDatabase(cc.Database).
			WithUsername(cc.User).
			WithPassword(cc.Password).
			WithDialTimeout(5 * time.Second).
			WithMaxOpenConn(10).
			WithMaxIdleConn(5).
			WithConnMaxLifetime(30 * time.Minute).
			WithBlockBufferSize(10)

		return New(ctx, cfg)
	})
}

func New(ctx context.Context, cfg *Config) (*Clickhouse, error) {
	opt := &clickhouse.Options{
		Addr: cfg.addresses,
		Auth: struct {
			Database string
			Username string
			Password string
		}{Database: cfg.database, Username: cfg.username, Password: cfg.password},
		DialTimeout:     cfg.dialTimeout,
		MaxOpenConns:    cfg.maxOpenConn,
		MaxIdleConns:    cfg.maxIdleConn,
		ConnMaxLifetime: cfg.connMaxLifetime,
		BlockBufferSize: uint8(cfg.blockBufferSize),
		TLS:             cfg.tlsConfig,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	}

	conn, err := clickhouse.Open(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to open clickhouse connection: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
	}

	return &Clickhouse{
		conn:   conn,
		config: cfg,
	}, nil
}

func (c *Clickhouse) Create(ctx context.Context) error {
	const createTableQuery = `CREATE TABLE %s (
    Id          Int32,
    CodeBlock   String,
    FilePath    String,
    Line        Int32,
    Description String
) ENGINE = MergeTree()
ORDER BY Id;`

	query := fmt.Sprintf(createTableQuery, c.config.tableName)
	if err := c.conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (c *Clickhouse) Save(ctx context.Context, id int, codeBlock string, filePath string, line int, description string) error {
	const insertQuery = `INSERT INTO %s (Id, CodeBlock, FilePath, Line, Description) VALUES (?, ?, ?, ?, ?);`

	query := fmt.Sprintf(insertQuery, c.config.tableName)
	if err := c.conn.Exec(ctx, query, id, codeBlock, filePath, line, description); err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

func (c *Clickhouse) Get(ctx context.Context, id int) (*code.Data, error) {
	const selectQuery = `SELECT Id, CodeBlock, FilePath, Line, Description FROM %s WHERE Id = ?;`

	query := fmt.Sprintf(selectQuery, c.config.tableName)
	rows, err := c.conn.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	data := &code.Data{}
	if rows.Next() {
		if err := rows.Scan(&data.Id, &data.CodeBlock, &data.FilePath, &data.Line, &data.Description); err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}
	}

	return data, nil
}

func (c *Clickhouse) Gets(ctx context.Context, ids ...int) ([]*code.Data, error) {
	const selectQuery = `SELECT Id, CodeBlock, FilePath, Line, Description FROM %s WHERE Id IN ?;`

	query := fmt.Sprintf(selectQuery, c.config.tableName)
	rows, err := c.conn.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var datas []*code.Data
	for rows.Next() {
		data := &code.Data{}
		if err := rows.Scan(&data.Id, &data.CodeBlock, &data.FilePath, &data.Line, &data.Description); err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (c *Clickhouse) Delete(ctx context.Context, id int) error {
	const deleteQuery = `DELETE FROM %s WHERE Id = ?;`

	query := fmt.Sprintf(deleteQuery, c.config.tableName)
	if err := c.conn.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (c *Clickhouse) Deletes(ctx context.Context, ids ...int) error {
	const deleteQuery = `DELETE FROM %s WHERE Id IN ?;`

	query := fmt.Sprintf(deleteQuery, c.config.tableName)
	if err := c.conn.Exec(ctx, query, ids); err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	return nil
}

func (c *Clickhouse) Drop(ctx context.Context) error {
	const dropTableQuery = `DROP TABLE %s;`

	query := fmt.Sprintf(dropTableQuery, c.config.tableName)
	if err := c.conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	return nil
}
