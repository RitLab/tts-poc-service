package command

import (
	"context"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/common/decorator"
)

type ReloadConfigHandler decorator.CommandHandler[*config.Cfg]

type reloadConfigRepository struct {
	log *baselogger.Logger
}

func NewReloadConfigRepository(log *baselogger.Logger) decorator.CommandHandler[*config.Cfg] {
	return decorator.ApplyCommandDecorators[*config.Cfg](
		reloadConfigRepository{log: log},
		log)
}

func (r reloadConfigRepository) Handle(ctx context.Context, _ *config.Cfg) (err error) {
	r.log.Logger.Info("start reload config")
	config.Config.Reload(ctx, r.log)
	r.log.Logger.Info("success reload config")
	return
}
