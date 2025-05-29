package query

import (
	"context"
	"fmt"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/decorator"
)

type VideoTranscriptQuery struct {
	Url string `json:"url" validate:"required"`
}

type VideoTranscriptResponse struct {
	Output string `json:"output"`
}

type VideoTranscriptHandler decorator.QueryHandler[VideoTranscriptQuery, VideoTranscriptResponse]

type videoTranscriptRepository struct {
	ai     gemini_ai.GenAIMethod
	logger *baselogger.Logger
}

func NewVideoTranscriptRepository(ai gemini_ai.GenAIMethod, log *baselogger.Logger) decorator.QueryHandler[VideoTranscriptQuery, VideoTranscriptResponse] {
	return decorator.ApplyQueryDecorators[VideoTranscriptQuery, VideoTranscriptResponse](
		videoTranscriptRepository{ai: ai, logger: log},
		log)
}

func (g videoTranscriptRepository) Handle(ctx context.Context, in VideoTranscriptQuery) (VideoTranscriptResponse, error) {
	result, err := g.ai.SummarizeYoutubeUrl(ctx, in.Url)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize youtube url: %w", err))
		return VideoTranscriptResponse{}, err
	}

	return VideoTranscriptResponse{Output: result}, nil
}
