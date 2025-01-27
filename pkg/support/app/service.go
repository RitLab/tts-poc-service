package app

import (
	"database/sql"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/support/adapter"
	"tts-poc-service/pkg/support/app/command"
)

type SupportService struct {
	Commands SupportCommands
}

type SupportCommands struct {
	InsertSupportHandler command.InsertSupportHandler
}

func NewSupportService(log *baselogger.Logger, db *sql.DB) SupportService {
	repo := adapter.NewSupportRepository(db, log)
	return SupportService{
		Commands: SupportCommands{
			InsertSupportHandler: command.NewInsertSupportRepository(repo, log),
		},
	}
}
