package query

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/tts/domain"
)

type GetTextToSpeechFileQuery struct {
	Text string
	Lang string
}

type GetTextToSpeechFileResponse struct {
	Url string `json:"url"`
}

type GetTextToSpeechHandler decorator.QueryHandler[GetTextToSpeechFileQuery, GetTextToSpeechFileResponse]

type getTextToSpeechRepository struct {
	s3     storage.Storage
	player htgo.Player
	logger *baselogger.Logger
}

func NewGetTextToSpeechRepository(s3 storage.Storage, player htgo.Player, log *baselogger.Logger) decorator.QueryHandler[GetTextToSpeechFileQuery, GetTextToSpeechFileResponse] {
	return decorator.ApplyQueryDecorators[GetTextToSpeechFileQuery, GetTextToSpeechFileResponse](
		getTextToSpeechRepository{s3: s3, player: player, logger: log},
		log)
}

func (g getTextToSpeechRepository) Handle(ctx context.Context, in GetTextToSpeechFileQuery) (GetTextToSpeechFileResponse, error) {
	filePaths, err := g.player.Save(in.Text, in.Lang)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error save file: %w", err))
		return GetTextToSpeechFileResponse{}, err
	}

	outputFile := fmt.Sprintf("%s/output-%s.mp3", constant.AUDIO_FOLDER, uuid.NewString())

	// Open the output file
	out, err := os.Create(outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error creating file: %w", err))
		return GetTextToSpeechFileResponse{}, nil
	}
	defer out.Close()

	// Loop through the input files and concatenate their data
	for _, path := range filePaths {
		err = domain.AppendFile(out, path)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error appending file: %w", err))
			return GetTextToSpeechFileResponse{}, nil
		}
	}

	err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: outputFile, ContentType: "audio/mpeg"})
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error put file: %w", err))
		return GetTextToSpeechFileResponse{}, err
	}

	filePaths = append(filePaths, outputFile)
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error removing file: %w", err))
		}
	}

	return GetTextToSpeechFileResponse{
		Url: fmt.Sprintf("%s://%s/%s/%s", config.Config.Storage.Method, config.Config.Storage.Endpoint, config.Config.Storage.BucketName, outputFile)}, nil
}
