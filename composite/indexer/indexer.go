package indexer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/snowmerak/indexer/lib/analyzer"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/lib/index/vector"
	"github.com/snowmerak/indexer/lib/store/code"
	"github.com/snowmerak/indexer/pkg/prompt"
	"github.com/snowmerak/indexer/pkg/stepper"
)

type Indexer struct {
	analyzer             analyzer.Analyzer
	embeddingsGeneration generation.Embeddings
	chatGeneration       generation.Text
	codeStore            code.Store
	vectorIndex          vector.Vector
}

func New(analyzer analyzer.Analyzer, embeddingsGeneration generation.Embeddings, chatGeneration generation.Text, codeStore code.Store, vectorIndex vector.Vector) *Indexer {
	return &Indexer{
		analyzer:             analyzer,
		embeddingsGeneration: embeddingsGeneration,
		chatGeneration:       chatGeneration,
		codeStore:            codeStore,
		vectorIndex:          vectorIndex,
	}
}

func (idx *Indexer) Initialize(ctx context.Context) error {
	rollback := false

	if err := idx.codeStore.Create(ctx); err != nil {
		rollback = true
		return fmt.Errorf("failed to create codeStore: %w", err)
	}
	defer func() {
		if rollback {
			if err := idx.codeStore.Drop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to drop codeStore")
			}
		}
	}()

	if err := idx.vectorIndex.Create(ctx, idx.embeddingsGeneration.Size()); err != nil {
		rollback = true
		return fmt.Errorf("failed to create vectorIndex: %w", err)
	}
	defer func() {
		if rollback {
			if err := idx.vectorIndex.Drop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to drop vectorIndex")
			}
		}
	}()

	return nil
}

func (idx *Indexer) CleanUp(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return idx.codeStore.Drop(ctx)
	})

	eg.Go(func() error {
		return idx.vectorIndex.Drop(ctx)
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to clean up: %w", err)
	}

	return nil
}

func (idx *Indexer) Index(ctx context.Context, path string) error {
	is := stepper.Int[int]()
	vis := stepper.Int[int]()

	if err := idx.analyzer.Walk(path, true, func(codeBlock, filePath string, line int) error {
		idxNum := is.Next()

		explanation, err := idx.chatGeneration.Generate(ctx, prompt.CodeAnalysis(codeBlock))
		if err != nil {
			return fmt.Errorf("failed to generate explanation: %w", err)
		}

		eg := errgroup.Group{}

		eg.Go(func() error {
			if err := idx.codeStore.Save(ctx, idxNum, codeBlock, filePath, line, explanation); err != nil {
				return fmt.Errorf("failed to save code: %w", err)
			}

			return nil
		})

		eg.Go(func() error {
			vectorIdxNum := vis.Next()

			embedding, err := idx.embeddingsGeneration.Embed(ctx, explanation)
			if err != nil {
				return fmt.Errorf("failed to embed explanation: %w", err)
			}

			if err := idx.vectorIndex.Store(ctx, vectorIdxNum, &vector.Payload{
				Id:        vectorIdxNum,
				Vector:    embedding,
				RelatedId: idxNum,
			}); err != nil {
				return fmt.Errorf("failed to store vector: %w", err)
			}

			return nil
		})

		eg.Go(func() error {
			vectorIdxNum := vis.Next()

			embedding, err := idx.embeddingsGeneration.Embed(ctx, codeBlock)
			if err != nil {
				return fmt.Errorf("failed to embed code: %w", err)
			}

			if err := idx.vectorIndex.Store(ctx, vectorIdxNum, &vector.Payload{
				Id:        vectorIdxNum,
				Vector:    embedding,
				RelatedId: idxNum,
			}); err != nil {
				return fmt.Errorf("failed to store vector: %w", err)
			}

			return nil
		})

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("failed to index: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk: %w", err)
	}

	return nil
}

func (idx *Indexer) Search(ctx context.Context, query string, count int) ([]*code.Data, error) {
	embedding, err := idx.embeddingsGeneration.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	ids, err := idx.vectorIndex.Search(ctx, embedding, count)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	data := make([]*code.Data, 0, len(ids))

	for _, id := range ids {
		codeData, err := idx.codeStore.Get(ctx, id.RelatedId)
		if err != nil {
			return nil, fmt.Errorf("failed to get code: %w", err)
		}

		data = append(data, codeData)
	}

	uniqueIdxSet := make(map[int]struct{})
	uniqueIdxOrder := make([]int, 0)
	for _, d := range data {
		if _, ok := uniqueIdxSet[d.Id]; !ok {
			uniqueIdxSet[d.Id] = struct{}{}
			uniqueIdxOrder = append(uniqueIdxOrder, d.Id)
		}
	}

	uniqueData := make([]*code.Data, 0, len(uniqueIdxOrder))
	for _, idx := range uniqueIdxOrder {
		uniqueData = append(uniqueData, data[idx])
	}

	return uniqueData, nil
}
