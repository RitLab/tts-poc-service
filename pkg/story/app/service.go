package app

import (
	"database/sql"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/story/adapters"
	"tts-poc-service/pkg/story/app/query"
)

type StoryService struct {
	Queries  StoryQueries
	Commands StoryCommands
}

type StoryQueries struct {
	GenerateStoryBookHandler query.GenerateStoryBookHandler
	ListStoryBookHandler     query.ListStoryBookHandler
	GetStoryBookHandler      query.GetStoryBookHandler
}

type StoryCommands struct{}

func NewStoryService(log *baselogger.Logger, s3 storage.Storage, db *sql.DB, ai gemini_ai.GenAIMethod, player htgo.Player) StoryService {
	storyAdapter := adapters.NewStoryRepository(db, log)
	return StoryService{
		Queries: StoryQueries{
			GenerateStoryBookHandler: query.NewGenerateStoryBookRepository(ai, player, storyAdapter, s3, log),
			ListStoryBookHandler:     query.NewListStoryBookRepository(storyAdapter, s3, log),
			GetStoryBookHandler:      query.NewGetStoryBookRepository(storyAdapter, s3, log),
		},
		Commands: StoryCommands{},
	}
}
