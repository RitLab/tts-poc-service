package query

import (
	"context"
	"fmt"
	"mime/multipart"
	"rsc.io/pdf"
	"strings"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
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
	openAIRepo domain.OpenAIRepository
	s3         storage.Storage
	logger     *baselogger.Logger
}

func NewSummarizePdfRepository(openAIRepo domain.OpenAIRepository, s3 storage.Storage,
	log *baselogger.Logger) decorator.QueryHandler[SummarizePdfQuery, SummarizePdfResponse] {
	return decorator.ApplyQueryDecorators[SummarizePdfQuery, SummarizePdfResponse](
		summarizePdfRepository{openAIRepo: openAIRepo, s3: s3, logger: log},
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

	doc, err := pdf.NewReader(src, in.File.Size)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return SummarizePdfResponse{}, err
	}

	content := make([]string, 0)
	numPages := doc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := doc.Page(i)
		if page.V.IsNull() {
			continue
		}
		for _, v := range page.Content().Text {
			content = append(content, v.S)
		}
	}

	openAIRequest := domain.InitOpenAIRequest(strings.Join(content, " "))
	summarize, err := g.openAIRepo.SummarizeText(ctx, openAIRequest)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize text: %w", err))
		return SummarizePdfResponse{}, err
	}

	return SummarizePdfResponse{Output: summarize}, nil
}
