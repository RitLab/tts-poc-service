package query

import (
	"context"
	"fmt"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/decorator"
)

type VideoSummarizeQuery struct {
	Url string `json:"url" validate:"required"`
}

type VideoSummarizeResponse struct {
	Output string `json:"output"`
}

type VideoSummarizeHandler decorator.QueryHandler[VideoSummarizeQuery, VideoSummarizeResponse]

type videoSummarizeRepository struct {
	ai     gemini_ai.GenAIMethod
	logger *baselogger.Logger
}

func NewVideoSummarizeRepository(ai gemini_ai.GenAIMethod, log *baselogger.Logger) decorator.QueryHandler[VideoSummarizeQuery, VideoSummarizeResponse] {
	return decorator.ApplyQueryDecorators[VideoSummarizeQuery, VideoSummarizeResponse](
		videoSummarizeRepository{ai: ai, logger: log},
		log)
}

func (g videoSummarizeRepository) Handle(ctx context.Context, in VideoSummarizeQuery) (VideoSummarizeResponse, error) {
	result, err := g.ai.SummarizeYoutubeUrl(ctx, in.Url)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize youtube url: %w", err))
		return VideoSummarizeResponse{}, err
	}

	return VideoSummarizeResponse{Output: result}, nil
}
