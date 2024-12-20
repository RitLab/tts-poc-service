package query

import (
	"context"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/common/decorator"
)

type GetConfigHandler decorator.QueryHandler[*config.Cfg, *config.Cfg]

type getConfigRepository struct {
	log *baselogger.Logger
}

func NewGetConfigRepository(log *baselogger.Logger) decorator.QueryHandler[*config.Cfg, *config.Cfg] {
	return decorator.ApplyQueryDecorators[*config.Cfg, *config.Cfg](
		getConfigRepository{log: log},
		log)
}

func (g getConfigRepository) Handle(ctx context.Context, _ *config.Cfg) (*config.Cfg, error) {
	return config.Config, nil
}
