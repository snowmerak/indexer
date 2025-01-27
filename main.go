package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/snowmerak/indexer/composite/indexer"
	"github.com/snowmerak/indexer/pkg/client/golang"
	"github.com/snowmerak/indexer/pkg/client/ollama"
	"github.com/snowmerak/indexer/pkg/client/postgres"
	"github.com/snowmerak/indexer/pkg/client/qdrant"
	"github.com/snowmerak/indexer/pkg/utill/jobs"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	jq, err := jobs.New(ctx, 36)
	if err != nil {
		log.Fatalf("failed to create jobs queue: %v", err)
	}

	otc, err := ollama.NewTextClient(ctx, ollama.NewClientConfig(), ollama.GenerationModelQwen2o5Coder1o5B)
	if err != nil {
		log.Fatalf("failed to create text client: %v", err)
	}

	oec, err := ollama.NewEmbeddingsClient(ctx, ollama.NewClientConfig(), ollama.EmbeddingModelBgeM3o5B)
	if err != nil {
		log.Fatalf("failed to create embeddings client: %v", err)
	}

	tableName := "indexer"

	vdb, err := qdrant.New(ctx, qdrant.NewConfig("localhost", 6334, tableName))
	if err != nil {
		log.Fatalf("failed to create vector database: %v", err)
	}

	pg, err := postgres.New(ctx, postgres.NewConfig("postgres://postgres:postgres@localhost:5432/postgres", tableName))
	if err != nil {
		log.Fatalf("failed to create postgres store: %v", err)
	}

	gaz := new(golang.Analyzer)

	idxer := indexer.New(jq, gaz, oec, otc, pg, vdb)

	switch command {
	case "init":
		if err := idxer.Initialize(ctx); err != nil {
			log.Fatalf("failed to initialize indexer: %v", err)
		}
	case "index":
		if firstArg == "" {
			log.Fatalf("index command requires a path to index")
		}

		if err := idxer.Index(ctx, firstArg); err != nil {
			panic(err)
		}
	case "search":
		if firstArg == "" {
			log.Fatalf("search command requires a query")
		}

		if secondArg == "" {
			secondArg = "10"
		}

		limitation, err := strconv.Atoi(secondArg)
		if err != nil {
			log.Fatalf("failed to parse limit: %v", err)
		}

		result, err := idxer.Search(ctx, firstArg, limitation)
		if err != nil {
			panic(err)
		}

		for _, r := range result {
			fmt.Println(r.FilePath, r.Line)
			fmt.Printf("-----------\nCode Block: %s\n-----------\n", r.CodeBlock)
			fmt.Printf("-----------\nExplanation: %s\n-----------\n", r.Description)
			fmt.Printf("==========\n")
		}
	case "cleanup":
		if err := idxer.CleanUp(ctx); err != nil {
			log.Fatalf("failed to clean up indexer: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s", command)
	}
}
