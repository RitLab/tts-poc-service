package app

import (
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/config/app/command"
	"tts-poc-service/pkg/config/app/query"
)

type ConfigService struct {
	Commands ConfigCommands
	Queries  ConfigQueries
}

type ConfigCommands struct {
	ReloadConfigHandler command.ReloadConfigHandler
}

type ConfigQueries struct {
	GetConfigHandler query.GetConfigHandler
}

func NewConfigService(log *baselogger.Logger) ConfigService {
	return ConfigService{
		Commands: ConfigCommands{
			ReloadConfigHandler: command.NewReloadConfigRepository(log),
		},
		Queries: ConfigQueries{
			GetConfigHandler: query.NewGetConfigRepository(log),
		},
	}
}
