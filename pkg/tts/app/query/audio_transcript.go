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

type AudioTranscriptQuery struct {
	File *multipart.FileHeader
}

type AudioTranscriptResponse struct {
	Output string `json:"output"`
}

type AudioTranscriptHandler decorator.QueryHandler[AudioTranscriptQuery, AudioTranscriptResponse]

type audioTranscriptRepository struct {
	ai     gemini_ai.GenAIMethod
	logger *baselogger.Logger
}

func NewAudioTranscriptRepository(ai gemini_ai.GenAIMethod, log *baselogger.Logger) decorator.QueryHandler[AudioTranscriptQuery, AudioTranscriptResponse] {
	return decorator.ApplyQueryDecorators[AudioTranscriptQuery, AudioTranscriptResponse](
		audioTranscriptRepository{ai: ai, logger: log},
		log)
}

func (g audioTranscriptRepository) Handle(ctx context.Context, in AudioTranscriptQuery) (AudioTranscriptResponse, error) {
	if err := domain.ValidateAudioFile(in.File); err != nil {
		return AudioTranscriptResponse{}, err
	}

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return AudioTranscriptResponse{}, err
	}
	defer src.Close()

	outputFile := fmt.Sprintf("%s/%s-%s", constant.AUDIO_FOLDER, uuid.NewString(), in.File.Filename)
	dst, err := os.Create(outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error create file: %w", err))
		return AudioTranscriptResponse{}, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error copy to a new file: %w", err))
		return AudioTranscriptResponse{}, err
	}

	result, err := g.ai.GetTranscriptAudio(ctx, outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error get transcript audio: %w", err))
		return AudioTranscriptResponse{}, err
	}

	return AudioTranscriptResponse{Output: result}, nil
}
