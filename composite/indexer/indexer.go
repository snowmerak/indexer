package indexer

import (
	"github.com/snowmerak/indexer/lib/analyzer"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/lib/index/vector"
	"github.com/snowmerak/indexer/lib/store/code"
)

type Indexer struct {
	analyzer             analyzer.Analyzer
	codeStore            code.Store
	embeddingsGeneration generation.Embeddings
	chatGeneration       generation.Text
	vectorIndex          vector.Vector
}

func New(analyzer analyzer.Analyzer, codeStore code.Store, embeddingsGeneration generation.Embeddings, chatGeneration generation.Text, vectorIndex vector.Vector) *Indexer {
	return &Indexer{
		analyzer:             analyzer,
		codeStore:            codeStore,
		embeddingsGeneration: embeddingsGeneration,
		chatGeneration:       chatGeneration,
		vectorIndex:          vectorIndex,
	}
}
