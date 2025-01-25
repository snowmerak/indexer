package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/snowmerak/indexer/composite/indexer"
	"github.com/snowmerak/indexer/lib/analyzer/golang"
	"github.com/snowmerak/indexer/lib/generation/ollama"
	"github.com/snowmerak/indexer/lib/index/vector/qdrant"
	"github.com/snowmerak/indexer/lib/store/code/postgres"
	"github.com/snowmerak/indexer/pkg/jobs"
)

func main() {
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

	oec, err := ollama.NewEmbeddingsClient(ctx, ollama.NewClientConfig(), ollama.EmbeddingModelMxbaiEmbedLarge)
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

	//if err := idxer.Initialize(ctx); err != nil {
	//	log.Fatalf("failed to initialize indexer: %v", err)
	//}
	//
	//if err := idxer.Index(ctx, "."); err != nil {
	//	panic(err)
	//}

	result, err := idxer.Search(ctx, "code explanation prompt", 10)
	if err != nil {
		panic(err)
	}

	for _, r := range result {
		fmt.Println(r.FilePath, r.Line)
		fmt.Printf("-----------\nCode Block: %s\n-----------\n", r.CodeBlock)
		fmt.Printf("-----------\nExplanation: %s\n-----------\n", r.Description)
		fmt.Printf("==========\n")
	}
}
