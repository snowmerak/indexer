package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/snowmerak/indexer/cps/indexer"
	"github.com/snowmerak/indexer/pkg/client/golang"
	"github.com/snowmerak/indexer/pkg/client/meilisearch"
	"github.com/snowmerak/indexer/pkg/client/ollama"
	"github.com/snowmerak/indexer/pkg/client/postgres"
	"github.com/snowmerak/indexer/pkg/client/pyembeddings"
	"github.com/snowmerak/indexer/pkg/client/qdrant"
	"github.com/snowmerak/indexer/pkg/util/jobs"
	"github.com/snowmerak/indexer/pkg/util/logger"
)

func main() {
	command := os.Args[1]
	firstArg := ""
	if len(os.Args) > 2 {
		firstArg = os.Args[2]
	}
	secondArg := ""
	if len(os.Args) > 3 {
		secondArg = os.Args[3]
	}

	logger.Init(zerolog.InfoLevel)

	tableName, err := filepath.Abs(firstArg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get absolute path")
	}

	{
		split := strings.Split(filepath.Dir(tableName), string(os.PathSeparator))
		if len(split) > 0 {
			tableName = split[len(split)-1]
		}
		if tableName == "" {
			tableName = "root_directory"
		}
	}

	log.Info().Str("detected_project_name", tableName).Msg("start application")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	jq, err := jobs.New(ctx, 36)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create job queue")
	}

	otc, err := ollama.NewTextClient(ctx, ollama.NewClientConfig(), ollama.GenerationModelQwen2o5Coder1o5B)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create text client")
	}

	ocec, err := pyembeddings.NewEmbeddings(ctx, pyembeddings.NewConfig("http://localhost:8392"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create code embeddings client")
	}

	otec, err := ollama.NewEmbeddingsClient(ctx, ollama.NewClientConfig(), ollama.EmbeddingModelBgeM3o5B)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create text embeddings client")
	}

	cvdb, err := qdrant.New(ctx, qdrant.NewConfig("localhost", 6334, tableName+"_code"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create vector database")
	}

	tvdb, err := qdrant.New(ctx, qdrant.NewConfig("localhost", 6334, tableName+"_desc"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create vector database")
	}

	pg, err := postgres.New(ctx, postgres.NewConfig("postgres://postgres:postgres@localhost:5432/postgres", tableName))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create postgres client")
	}

	ms, err := meilisearch.New(ctx, meilisearch.NewConfig("http://localhost:7700", "indexer").WithApiKey("tFWSre9Ix9Ltq7nXV87c9O5UP"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create meilisearch client")
	}

	gaz := new(golang.Analyzer)

	idxer := indexer.New(jq, gaz, ocec, otec, otc, pg, cvdb, tvdb, ms)

	switch command {
	case "init":
		if err := idxer.Initialize(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to initialize indexer")
		}
	case "index":
		if firstArg == "" {
			log.Fatal().Msg("index command requires a path")
		}

		if err := idxer.Index(ctx, firstArg); err != nil {
			panic(err)
		}
	case "search":
		if firstArg == "" {
			log.Fatal().Msg("search command requires a query")
		}

		if secondArg == "" {
			secondArg = "10"
		}

		limitation, err := strconv.Atoi(secondArg)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse limitation")
		}

		result, err := idxer.Search(ctx, firstArg, limitation)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to search")
		}

		for _, r := range result {
			fmt.Println(r.FilePath, r.Line)
			fmt.Printf("-----------\nCode Block: %s\n-----------\n", r.CodeBlock)
			fmt.Printf("-----------\nExplanation: %s\n-----------\n", r.Description)
			fmt.Printf("==========\n")
		}
	case "cleanup":
		if err := idxer.CleanUp(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to cleanup")
		}
	default:
		log.Fatal().Str("input", command).Msg("unknown command")
	}
}
