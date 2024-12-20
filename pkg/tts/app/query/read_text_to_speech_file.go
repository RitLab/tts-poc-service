package query

import (
	"context"
	"fmt"
	"os"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
)

type ReadTextToSpeechFileQuery struct {
	Text string
	Lang string
}

type ReadTextToSpeechFileResponse struct {
}

type ReadTextToSpeechHandler decorator.QueryHandler[ReadTextToSpeechFileQuery, ReadTextToSpeechFileResponse]

type readTextToSpeechRepository struct {
	player htgo.Player
	logger *baselogger.Logger
}

func NewReadTextToSpeechRepository(player htgo.Player, log *baselogger.Logger) decorator.QueryHandler[ReadTextToSpeechFileQuery, ReadTextToSpeechFileResponse] {
	return decorator.ApplyQueryDecorators[ReadTextToSpeechFileQuery, ReadTextToSpeechFileResponse](
		readTextToSpeechRepository{player: player, logger: log},
		log)
}

func (g readTextToSpeechRepository) Handle(ctx context.Context, in ReadTextToSpeechFileQuery) (ReadTextToSpeechFileResponse, error) {
	err := g.player.Play(in.Text, in.Lang)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error play text to speech file: %w", err))
		return ReadTextToSpeechFileResponse{}, err
	}

	files, err := os.ReadDir(constant.AUDIO_FOLDER)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read audio folder: %w", err))
		return ReadTextToSpeechFileResponse{}, nil
	}

	// Loop through files and delete each
	for _, file := range files {
		filePath := constant.AUDIO_FOLDER + "/" + file.Name()
		err = os.Remove(filePath)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error remove audio file: %w", err))
		}
	}

	return ReadTextToSpeechFileResponse{}, nil
}
