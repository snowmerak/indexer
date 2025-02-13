package indexer

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/snowmerak/indexer/lib/analyzer"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/lib/index/text"
	"github.com/snowmerak/indexer/lib/index/vector"
	"github.com/snowmerak/indexer/lib/store/code"
	"github.com/snowmerak/indexer/pkg/util/ext"
	"github.com/snowmerak/indexer/pkg/util/jobs"
	"github.com/snowmerak/indexer/pkg/util/prompt"
	"github.com/snowmerak/indexer/pkg/util/stepper"
)

type Indexer struct {
	jobs                     *jobs.Jobs
	analyzer                 analyzer.Analyzer
	codeEmbeddingsGeneration generation.Embeddings
	textEmbeddingsGeneration generation.Embeddings
	chatGeneration           generation.Text
	codeStore                code.Store
	codeVectorIndex          vector.Vector
	textVectorIndex          vector.Vector
	textIndex                text.Text
}

func New(jobs *jobs.Jobs, analyzer analyzer.Analyzer, codeEmbeddingsGeneration generation.Embeddings, textEmbeddingsGeneration generation.Embeddings, chatGeneration generation.Text, codeStore code.Store, codeVectorIndex vector.Vector, textVectorIndex vector.Vector, textIndex text.Text) *Indexer {
	return &Indexer{
		jobs:                     jobs,
		analyzer:                 analyzer,
		codeEmbeddingsGeneration: codeEmbeddingsGeneration,
		textEmbeddingsGeneration: textEmbeddingsGeneration,
		chatGeneration:           chatGeneration,
		codeStore:                codeStore,
		codeVectorIndex:          codeVectorIndex,
		textVectorIndex:          textVectorIndex,
		textIndex:                textIndex,
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

	cves, err := idx.codeEmbeddingsGeneration.Size()
	if err != nil {
		rollback = true
		return fmt.Errorf("failed to get codeEmbeddingsGeneration size: %w", err)
	}

	if err := idx.codeVectorIndex.Create(ctx, cves); err != nil {
		rollback = true
		return fmt.Errorf("failed to create vectorIndex: %w", err)
	}
	defer func() {
		if rollback {
			if err := idx.codeVectorIndex.Drop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to drop vectorIndex")
			}
		}
	}()

	tves, err := idx.textEmbeddingsGeneration.Size()
	if err != nil {
		rollback = true
		return fmt.Errorf("failed to get textEmbeddingsGeneration size: %w", err)
	}

	if err := idx.textVectorIndex.Create(ctx, tves); err != nil {
		rollback = true
		return fmt.Errorf("failed to create vectorIndex: %w", err)
	}
	defer func() {
		if rollback {
			if err := idx.textVectorIndex.Drop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to drop vectorIndex")
			}
		}
	}()

	if err := idx.textIndex.Create(ctx); err != nil {
		rollback = true
		return fmt.Errorf("failed to create textIndex: %w", err)
	}
	defer func() {
		if rollback {
			if err := idx.textIndex.Drop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to drop textIndex")
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
		return idx.codeVectorIndex.Drop(ctx)
	})

	eg.Go(func() error {
		return idx.textVectorIndex.Drop(ctx)
	})

	eg.Go(func() error {
		return idx.textIndex.Drop(ctx)
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to clean up: %w", err)
	}

	return nil
}

func (idx *Indexer) Index(ctx context.Context, path string) error {
	is := stepper.Int[int]()

	if err := idx.analyzer.Walk(path, true, func(codeBlock, filePath string, line int) error {
		_ = idx.jobs.Submit(func() error {
			idxNum := is.Next()

			lang := ext.ToLang(filepath.Ext(filePath))

			explanation, err := idx.chatGeneration.Generate(ctx, prompt.CodeAnalysis(lang, codeBlock))
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
				embedding, err := idx.textEmbeddingsGeneration.Embed(ctx, explanation)
				if err != nil {
					return fmt.Errorf("failed to embed explanation: %w", err)
				}

				if err := idx.textVectorIndex.Store(ctx, idxNum, &vector.Payload{
					Id:        idxNum,
					Vector:    embedding,
					RelatedId: idxNum,
				}); err != nil {
					return fmt.Errorf("failed to store vector: %w", err)
				}

				return nil
			})

			eg.Go(func() error {
				embedding, err := idx.codeEmbeddingsGeneration.Embed(ctx, codeBlock)
				if err != nil {
					return fmt.Errorf("failed to embed code: %w", err)
				}

				if err := idx.codeVectorIndex.Store(ctx, idxNum, &vector.Payload{
					Id:        idxNum,
					Vector:    embedding,
					RelatedId: idxNum,
				}); err != nil {
					return fmt.Errorf("failed to store vector: %w", err)
				}

				return nil
			})

			eg.Go(func() error {
				if err := idx.textIndex.Store(ctx, idxNum, text.Payload{
					Description: explanation,
					CodeBlock:   codeBlock,
				}); err != nil {
					return fmt.Errorf("failed to store text: %w", err)
				}

				return nil
			})

			if err := eg.Wait(); err != nil {
				return fmt.Errorf("failed to index: %w", err)
			}

			return nil
		})

		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk: %w", err)
	}

	if err := idx.jobs.Close(); err != nil {
		return fmt.Errorf("failed to close jobs: %w", err)
	}

	return nil
}

func (idx *Indexer) Search(ctx context.Context, query string, count int) ([]*code.Data, error) {
	cvr := make([]*code.Data, 0)
	tvr := make([]*code.Data, 0)
	tr := make([]*code.Data, 0)

	eg := errgroup.Group{}

	eg.Go(func() error {
		err := error(nil)
		cvr, err = idx.searchVector(ctx, query, count, CodeVector)
		if err != nil {
			return fmt.Errorf("failed to search vector: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		err := error(nil)
		tvr, err = idx.searchVector(ctx, query, count, TextVector)
		if err != nil {
			return fmt.Errorf("failed to search vector: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		err := error(nil)
		tr, err = idx.searchText(ctx, query, count)
		if err != nil {
			return fmt.Errorf("failed to search text: %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	cr := make([]*code.Data, 0, len(cvr)+len(tr))
	inserted := make(map[int]struct{})
	for _, d := range cvr {
		if _, ok := inserted[d.Id]; !ok {
			cr = append(cr, d)
			inserted[d.Id] = struct{}{}
		}
	}

	for _, d := range tvr {
		if _, ok := inserted[d.Id]; !ok {
			cr = append(cr, d)
			inserted[d.Id] = struct{}{}
		}
	}

	for _, d := range tr {
		if _, ok := inserted[d.Id]; !ok {
			cr = append(cr, d)
			inserted[d.Id] = struct{}{}
		}
	}

	return cr, nil
}

type VectorType int

const (
	CodeVector VectorType = iota
	TextVector
)

func (idx *Indexer) searchVector(ctx context.Context, query string, count int, searchVectorType VectorType) ([]*code.Data, error) {
	embeddings := generation.Embeddings(nil)
	vectorIndex := vector.Vector(nil)

	switch searchVectorType {
	case CodeVector:
		embeddings = idx.codeEmbeddingsGeneration
		vectorIndex = idx.codeVectorIndex
	case TextVector:
		embeddings = idx.textEmbeddingsGeneration
		vectorIndex = idx.textVectorIndex
	default:
		return nil, fmt.Errorf("invalid search vector type: %d", searchVectorType)
	}

	embedding, err := embeddings.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	ids, err := vectorIndex.Search(ctx, embedding, count)
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
	for i, d := range data {
		if _, ok := uniqueIdxSet[d.Id]; !ok {
			uniqueIdxSet[d.Id] = struct{}{}
			uniqueIdxOrder = append(uniqueIdxOrder, i)
		}
	}

	uniqueData := make([]*code.Data, 0, len(uniqueIdxOrder))
	for _, idx := range uniqueIdxOrder {
		uniqueData = append(uniqueData, data[idx])
	}

	return uniqueData, nil
}

func (idx *Indexer) searchText(ctx context.Context, query string, count int) ([]*code.Data, error) {
	results, err := idx.textIndex.Query(ctx, query, text.SearchOption{
		Limit: count,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	data := make([]*code.Data, 0, len(results))

	for _, result := range results {
		codeData, err := idx.codeStore.Get(ctx, result.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get code: %w", err)
		}

		data = append(data, codeData)
	}

	return data, nil
}
