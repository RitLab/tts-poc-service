package decorator

import (
	"context"
	"tts-poc-service/lib/baselogger"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger *baselogger.Logger) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base: handler,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
