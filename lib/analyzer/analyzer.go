package analyzer

import (
	"fmt"
	"sync"

	"github.com/snowmerak/indexer/pkg/config"
)

type Analyzer interface {
	Walk(path string, recursive bool, callback func(codeBlock, filePath string, line int) error) error
}

var registered = sync.Map{}

type Constructor func(*config.ClientConfig) (Analyzer, error)

func Register(name string, constructor Constructor) {
	registered.Store(name, constructor)
}

func Get(name string, config *config.ClientConfig) (Analyzer, error) {
	if v, ok := registered.Load(name); ok {
		analyzer, err := v.(Constructor)(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create analyzer: %w", err)
		}
		return analyzer, nil
	}
	return nil, fmt.Errorf("analyzer not found: %s", name)
}
