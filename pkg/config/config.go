package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFilename = "config.yaml"
)

type ClientConfig struct {
	Type      string   `yaml:"type,omitempty"`
	Scheme    string   `yaml:"scheme,omitempty"`
	Host      []string `yaml:"host,omitempty"`
	Database  string   `yaml:"database,omitempty"`
	User      string   `yaml:"user,omitempty"`
	Password  string   `yaml:"password,omitempty"`
	ApiKey    string   `yaml:"api_key,omitempty"`
	Model     string   `yaml:"model,omitempty"`
	Dimension int      `yaml:"dimension,omitempty"`
	Project   string   `yaml:"project,omitempty"`
}

type Config struct {
	MaxConcurrentJobs int    `yaml:"max_concurrent_jobs,omitempty"`
	Analyzer          string `yaml:"analyzer,omitempty"`
	Embeddings        struct {
		Code        ClientConfig `yaml:"code,omitempty"`
		Description ClientConfig `yaml:"description,omitempty"`
	} `yaml:"embeddings,omitempty"`
	Generation struct {
		Chat ClientConfig `yaml:"chat,omitempty"`
	} `yaml:"generation,omitempty"`
	Store struct {
		Code ClientConfig `yaml:"code,omitempty"`
	}
	Index struct {
		Vector struct {
			Code        ClientConfig `yaml:"code,omitempty"`
			Description ClientConfig `yaml:"description,omitempty"`
		} `yaml:"vector,omitempty"`
		Text struct {
			Index ClientConfig `yaml:"index,omitempty"`
		} `yaml:"description,omitempty"`
	} `yaml:"index,omitempty"`
}

func Default() *Config {
	return &Config{
		MaxConcurrentJobs: 36,
		Analyzer:          "golang",
		Embeddings: struct {
			Code        ClientConfig `yaml:"code,omitempty"`
			Description ClientConfig `yaml:"description,omitempty"`
		}{
			Code: ClientConfig{
				Type: "pyembeddings",
				Host: []string{"http://localhost:8392"},
			},
			Description: ClientConfig{
				Type:      "ollama",
				Host:      []string{"http://localhost:11434"},
				Model:     "bge-m3",
				Dimension: 1024,
			},
		},
		Generation: struct {
			Chat ClientConfig `yaml:"chat,omitempty"`
		}{
			Chat: ClientConfig{
				Type:  "ollama",
				Host:  []string{"http://localhost:11434"},
				Model: "qwen2.5-coder:1.5b",
			},
		},
		Store: struct {
			Code ClientConfig `yaml:"code,omitempty"`
		}{
			Code: ClientConfig{
				Type:     "postgres",
				Host:     []string{"localhost:5432"},
				User:     "postgres",
				Password: "postgres",
				Database: "postgres",
				Project:  "your-project-name",
			},
		},
		Index: struct {
			Vector struct {
				Code        ClientConfig `yaml:"code,omitempty"`
				Description ClientConfig `yaml:"description,omitempty"`
			} `yaml:"vector,omitempty"`
			Text struct {
				Index ClientConfig `yaml:"index,omitempty"`
			} `yaml:"description,omitempty"`
		}{
			Vector: struct {
				Code        ClientConfig `yaml:"code,omitempty"`
				Description ClientConfig `yaml:"description,omitempty"`
			}{
				Code: ClientConfig{
					Type:    "qdrant",
					Host:    []string{"localhost:6334"},
					Project: "your-project-name-code",
				},
				Description: ClientConfig{
					Type:    "qdrant",
					Host:    []string{"localhost:6334"},
					Project: "your-project-name-description",
				},
			},
			Text: struct {
				Index ClientConfig `yaml:"index,omitempty"`
			}{
				Index: ClientConfig{
					Type:    "meilisearch",
					Host:    []string{"http://localhost:7700"},
					ApiKey:  "tFWSre9Ix9Ltq7nXV87c9O5UP",
					Project: "your-project-name",
				},
			},
		},
	}
}

func (c *Config) MarshalTo(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}

func (c *Config) MarshalToBytes() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *Config) UnmarshalFrom(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

func (c *Config) UnmarshalFromBytes(b []byte) error {
	return yaml.Unmarshal(b, c)
}
