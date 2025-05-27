package command

import (
	"context"
	"fmt"
	"mime/multipart"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/database"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/pdf/domain"
)

type UpdateContextPdfQuery struct {
	File *multipart.FileHeader
}

type UpdateContextPdfHandler decorator.CommandHandler[UpdateContextPdfQuery]

type updateContextPdfRepository struct {
	dbVector database.VectorDatabase
	ai       gemini_ai.GenAIMethod
	logger   *baselogger.Logger
}

func NewUpdateContextPdfRepository(dbVector database.VectorDatabase, ai gemini_ai.GenAIMethod, log *baselogger.Logger) decorator.CommandHandler[UpdateContextPdfQuery] {
	return decorator.ApplyCommandDecorators[UpdateContextPdfQuery](
		updateContextPdfRepository{dbVector: dbVector, ai: ai, logger: log},
		log)
}

func (g updateContextPdfRepository) Handle(ctx context.Context, in UpdateContextPdfQuery) (err error) {
	if err = domain.ValidatePdfFile(in.File); err != nil {
		return err
	}

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return err
	}
	defer src.Close()

	pdfReader, err := domain.NewPdfReader(src)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read file: %w", err))
		return err
	}
	cleanedText, err := pdfReader.CleanText()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error clean text: %w", err))
		return err
	}

	chunkSize := 1000
	sentenceOverlap := 1
	sentences := utils.SplitIntoSentences(cleanedText)
	textChunks := utils.ChunkSentences(sentences, chunkSize, sentenceOverlap)

	embeddings := make([][]float32, 0, len(textChunks))
	for _, chunk := range textChunks {
		if !utils.IsValidChunk(chunk) {
			continue
		}
		embedding, err := g.ai.TextEmbedding(ctx, chunk)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error get embedding: %w", err))
			return err
		}
		err = g.dbVector.CheckSimilarity(ctx, config.Config.General.MilvusCollectionName, embedding)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error check similarity: %w", err))
			continue
		}

		embeddings = append(embeddings, embedding)
	}

	if len(embeddings) == 0 {
		return nil
	}

	err = g.dbVector.StoreEmbedding(ctx, database.Embedding{
		Value: embeddings,
		Chunk: textChunks,
	}, config.Config.General.MilvusCollectionName)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error store embedding: %w", err))
		return err
	}

	return nil
}
