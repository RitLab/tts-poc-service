package app

import (
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
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
	JoinMp3FilesHandler     query.JoinMp3FilesHandler
	AudioTranscriptHandler  query.AudioTranscriptHandler
	AudioSummarizeHandler   query.AudioSummarizeHandler
	VideoTranscriptHandler  query.VideoTranscriptHandler
	VideoSummarizeHandler   query.VideoSummarizeHandler
}

func NewTtsService(log *baselogger.Logger, player htgo.Player, s3 storage.Storage, ai gemini_ai.GenAIMethod) TtService {
	return TtService{
		Queries: TtsQueries{
			ReadTextToSpeechHandler: query.NewReadTextToSpeechRepository(player, log),
			GetTextToSpeechHandler:  query.NewGetTextToSpeechRepository(s3, player, log),
			JoinMp3FilesHandler:     query.NewJoinMp3FilesRepository(s3, log),
			AudioTranscriptHandler:  query.NewAudioTranscriptRepository(ai, log),
			AudioSummarizeHandler:   query.NewAudioSummarizeRepository(ai, log),
			VideoTranscriptHandler:  query.NewVideoTranscriptRepository(ai, log),
			VideoSummarizeHandler:   query.NewVideoSummarizeRepository(ai, log),
		},
	}
}
