package query

import (
	"context"
	"fmt"
	"mime/multipart"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/pdf/domain"
)

type SummarizePdfQuery struct {
	File *multipart.FileHeader
}

type SummarizePdfResponse struct {
	Output string `json:"output"`
}

type SummarizePdfHandler decorator.QueryHandler[SummarizePdfQuery, SummarizePdfResponse]

type summarizePdfRepository struct {
	ai     gemini_ai.GenAIMethod
	logger *baselogger.Logger
}

func NewSummarizePdfRepository(ai gemini_ai.GenAIMethod,
	log *baselogger.Logger) decorator.QueryHandler[SummarizePdfQuery, SummarizePdfResponse] {
	return decorator.ApplyQueryDecorators[SummarizePdfQuery, SummarizePdfResponse](
		summarizePdfRepository{ai: ai, logger: log},
		log)
}

func (g summarizePdfRepository) Handle(ctx context.Context, in SummarizePdfQuery) (SummarizePdfResponse, error) {
	if err := domain.ValidatePdfFile(in.File); err != nil {
		return SummarizePdfResponse{}, err
	}

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return SummarizePdfResponse{}, err
	}
	defer src.Close()

	pdfReader, err := domain.NewPdfReader(src)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read file: %w", err))
		return SummarizePdfResponse{}, err
	}
	cleanedText, err := pdfReader.CleanText()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error clean text: %w", err))
		return SummarizePdfResponse{}, err
	}

	result, err := g.ai.SummarizeText(ctx, cleanedText)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize text: %w", err))
		return SummarizePdfResponse{}, err
	}

	return SummarizePdfResponse{Output: result}, nil
}
