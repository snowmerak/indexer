package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/snowmerak/indexer/cps/indexer"
	"github.com/snowmerak/indexer/lib/analyzer"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/lib/index/text"
	"github.com/snowmerak/indexer/lib/index/vector"
	"github.com/snowmerak/indexer/lib/store/code"
	"github.com/snowmerak/indexer/pkg/config"
	"github.com/snowmerak/indexer/pkg/util/ext"
	"github.com/snowmerak/indexer/pkg/util/jobs"
	"github.com/snowmerak/indexer/pkg/util/logger"

	_ "github.com/snowmerak/indexer/pkg/client/clickhouse"
	_ "github.com/snowmerak/indexer/pkg/client/gemini"
	_ "github.com/snowmerak/indexer/pkg/client/golang"
	_ "github.com/snowmerak/indexer/pkg/client/meilisearch"
	_ "github.com/snowmerak/indexer/pkg/client/ollama"
	_ "github.com/snowmerak/indexer/pkg/client/postgres"
	_ "github.com/snowmerak/indexer/pkg/client/pyembeddings"
	_ "github.com/snowmerak/indexer/pkg/client/qdrant"
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

	tableName, err := ext.GetDirectoryName(firstArg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get project name")
	}

	log.Info().Str("detected_project_name", tableName).Msg("start application")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	switch command {
	case "new":
		cfg := config.Default()
		f, err := os.Create(config.DefaultFilename)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create config file")
		}
		defer f.Close()

		if err := cfg.MarshalTo(f); err != nil {
			log.Fatal().Err(err).Msg("failed to write config")
		}

		log.Info().Str("filename", config.DefaultFilename).Msg("new config file created")
	default:
		cfg := config.Default()
		f, err := os.Open(config.DefaultFilename)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open config file")
		}

		if err := cfg.UnmarshalFrom(f); err != nil {
			log.Fatal().Err(err).Msg("failed to read config")
		}

		jq, err := jobs.New(ctx, cfg.MaxConcurrentJobs)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create job queue")
		}

		otc, err := generation.GetText(ctx, cfg.Generation.Chat.Type, &cfg.Generation.Chat)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create text client")
		}

		ocec, err := generation.GetEmbeddings(ctx, cfg.Embeddings.Code.Type, &cfg.Embeddings.Code)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create code embeddings client")
		}

		otec, err := generation.GetEmbeddings(ctx, cfg.Embeddings.Description.Type, &cfg.Embeddings.Description)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create text embeddings client")
		}

		cvdb, err := vector.GetVector(ctx, cfg.Index.Vector.Code.Type, &cfg.Index.Vector.Code)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create vector database")
		}

		tvdb, err := vector.GetVector(ctx, cfg.Index.Vector.Description.Type, &cfg.Index.Vector.Description)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create vector database")
		}

		pg, err := code.GetStore(ctx, cfg.Store.Code.Type, &cfg.Store.Code)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create postgres client")
		}

		ms, err := text.GetText(ctx, cfg.Index.Text.Index.Type, &cfg.Index.Text.Index)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create meilisearch client")
		}

		gaz, err := analyzer.Get(ctx, cfg.Analyzer, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create analyzer")
		}

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
				secondArg = strconv.FormatInt(int64(cfg.SearchCount), 10)
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
				fmt.Println(r.ToMarkdown(gaz.LanguageCode()))
			}
		case "cleanup":
			if err := idxer.CleanUp(ctx); err != nil {
				log.Fatal().Err(err).Msg("failed to cleanup")
			}
		default:
			log.Fatal().Str("input", command).Msg("unknown command")
		}
	}
}
