package query

import (
	"context"
	"fmt"
	"strings"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/database"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/common/utils"
)

type GenerateContextQuery struct {
	Question string `json:"question" validate:"required"`
}

type GenerateContextResponse struct {
	Output string `json:"output"`
}

type GenerateContextHandler decorator.QueryHandler[GenerateContextQuery, GenerateContextResponse]

type generateContextRepository struct {
	dbVector database.VectorDatabase
	ai       gemini_ai.GenAIMethod
	logger   *baselogger.Logger
}

func NewGenerateContextRepository(dbVector database.VectorDatabase, ai gemini_ai.GenAIMethod,
	log *baselogger.Logger) decorator.QueryHandler[GenerateContextQuery, GenerateContextResponse] {
	return decorator.ApplyQueryDecorators[GenerateContextQuery, GenerateContextResponse](
		generateContextRepository{dbVector: dbVector, ai: ai, logger: log},
		log)
}

func (g generateContextRepository) Handle(ctx context.Context, in GenerateContextQuery) (GenerateContextResponse, error) {
	chunkSize := 1000
	sentencesOverlap := 1
	sentences := utils.SplitIntoSentences(in.Question)
	textChunks := utils.ChunkSentences(sentences, chunkSize, sentencesOverlap)

	searchText := make([]string, 0)
	for _, chunk := range textChunks {
		embedding, err := g.ai.TextEmbedding(ctx, chunk)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error get embedding: %w", err))
			return GenerateContextResponse{}, err
		}

		contextText, err := g.dbVector.SearchEmbedding(ctx, config.Config.General.MilvusCollectionName, embedding, 3)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error searching embedding: %w", err))
			return GenerateContextResponse{}, nil
		}

		searchText = append(searchText, contextText)
	}

	result, err := g.ai.GenerateFromContext(ctx, strings.Join(searchText, " "), in.Question)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize text: %w", err))
		return GenerateContextResponse{}, err
	}

	return GenerateContextResponse{Output: result}, nil
}
