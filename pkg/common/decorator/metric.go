package decorator

import (
	"context"
)

type commandMetricsDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandMetricsDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	return d.base.Handle(ctx, cmd)
}

type queryMetricsDecorator[C any, R any] struct {
	base QueryHandler[C, R]
}

func (d queryMetricsDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	return d.base.Handle(ctx, query)
}
