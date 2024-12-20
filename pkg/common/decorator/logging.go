package decorator

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"tts-poc-service/lib/baselogger"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *baselogger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := d.logger.Hashcode(ctx).WithFields(logrus.Fields{
		"command":      handlerType,
		"command_body": fmt.Sprintf("%#v", cmd),
	})

	defer func() {
		if err == nil {
			logger.Info(fmt.Sprintf("Command executed successfully: %s", handlerType))
		} else {
			logger.Error(fmt.Errorf("failed execute the config %s: %w", handlerType, err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger *baselogger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerType := generateActionName(cmd)

	logger := d.logger.Hashcode(ctx).WithFields(logrus.Fields{
		"search":     handlerType,
		"query_body": fmt.Sprintf("%#v", cmd),
	})

	defer func() {
		if err == nil {
			logger.Info(fmt.Sprintf("Query executed successfully: %s", handlerType))
		} else {
			logger.Error(fmt.Errorf("failed execute the search %s: %w", handlerType, err))
		}
	}()

	return d.base.Handle(ctx, cmd)
}
