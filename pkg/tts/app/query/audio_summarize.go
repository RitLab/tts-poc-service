package query

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/tts/domain"
)

type AudioSummarizeQuery struct {
	File *multipart.FileHeader
}

type AudioSummarizeResponse struct {
	Output string `json:"output"`
}

type AudioSummarizeHandler decorator.QueryHandler[AudioSummarizeQuery, AudioSummarizeResponse]

type audioSummarizeRepository struct {
	ai     gemini_ai.GenAIMethod
	logger *baselogger.Logger
}

func NewAudioSummarizeRepository(ai gemini_ai.GenAIMethod, log *baselogger.Logger) decorator.QueryHandler[AudioSummarizeQuery, AudioSummarizeResponse] {
	return decorator.ApplyQueryDecorators[AudioSummarizeQuery, AudioSummarizeResponse](
		audioSummarizeRepository{ai: ai, logger: log},
		log)
}

func (g audioSummarizeRepository) Handle(ctx context.Context, in AudioSummarizeQuery) (AudioSummarizeResponse, error) {
	if err := domain.ValidateAudioFile(in.File); err != nil {
		return AudioSummarizeResponse{}, err
	}

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return AudioSummarizeResponse{}, err
	}
	defer src.Close()

	outputFile := fmt.Sprintf("%s/%s-%s", constant.AUDIO_FOLDER, uuid.NewString(), in.File.Filename)
	dst, err := os.Create(outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error create file: %w", err))
		return AudioSummarizeResponse{}, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error copy to a new file: %w", err))
		return AudioSummarizeResponse{}, err
	}

	trancript, err := g.ai.GetTranscriptAudio(ctx, outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error get transcript audio: %w", err))
		return AudioSummarizeResponse{}, err
	}

	result, err := g.ai.SummarizeText(ctx, trancript)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error summarize text: %w", err))
		return AudioSummarizeResponse{}, err
	}

	return AudioSummarizeResponse{Output: result}, nil
}
