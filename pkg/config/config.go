package config

type ClientConfig struct {
	Type      string `yaml:"type,omitempty"`
	Scheme    string `yaml:"scheme,omitempty"`
	Host      string `yaml:"host,omitempty"`
	Port      int    `yaml:"port,omitempty"`
	Database  string `yaml:"database,omitempty"`
	User      string `yaml:"user,omitempty"`
	Password  string `yaml:"password,omitempty"`
	ApiKey    string `yaml:"api_key,omitempty"`
	Model     string `yaml:"model,omitempty"`
	Dimension int    `yaml:"dimension,omitempty"`
	Project   string `yaml:"project,omitempty"`
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

func DefaultConfig() *Config {
	return &Config{
		MaxConcurrentJobs: 36,
		Analyzer:          "golang",
		Embeddings: struct {
			Code        ClientConfig `yaml:"code,omitempty"`
			Description ClientConfig `yaml:"description,omitempty"`
		}{
			Code: ClientConfig{
				Type: "pyembeddings",
				Host: "http://localhost:8392",
			},
			Description: ClientConfig{
				Type:      "ollama",
				Host:      "http://localhost:11434",
				Model:     "bge-m3",
				Dimension: 1024,
			},
		},
		Generation: struct {
			Chat ClientConfig `yaml:"chat,omitempty"`
		}{
			Chat: ClientConfig{
				Type:  "ollama",
				Host:  "http://localhost:11434",
				Model: "qwen2.5-coder:1.5b",
			},
		},
		Store: struct {
			Code ClientConfig `yaml:"code,omitempty"`
		}{
			Code: ClientConfig{
				Type:     "postgres",
				Host:     "localhost",
				Port:     5432,
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
					Host:    "localhost",
					Port:    6334,
					Project: "your-project-name-code",
				},
				Description: ClientConfig{
					Type:    "qdrant",
					Host:    "localhost",
					Port:    6334,
					Project: "your-project-name-description",
				},
			},
			Text: struct {
				Index ClientConfig `yaml:"index,omitempty"`
			}{
				Index: ClientConfig{
					Type:    "meilisearch",
					Host:    "http://localhost:7700",
					ApiKey:  "tFWSre9Ix9Ltq7nXV87c9O5UP",
					Project: "your-project-name",
				},
			},
		},
	}
}
