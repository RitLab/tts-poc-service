package app

import (
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/tts/app/query"
)

type TtService struct {
	Queries TtsQueries
}

type TtsQueries struct {
	ReadTextToSpeechHandler query.ReadTextToSpeechHandler
	GetTextToSpeechHandler  query.GetTextToSpeechHandler
}

func NewTtsService(log *baselogger.Logger, player htgo.Player, s3 storage.Storage) TtService {
	return TtService{
		Queries: TtsQueries{
			ReadTextToSpeechHandler: query.NewReadTextToSpeechRepository(player, log),
			GetTextToSpeechHandler:  query.NewGetTextToSpeechRepository(s3, player, log),
		},
	}
}
