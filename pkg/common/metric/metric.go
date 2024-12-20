package pkg_metric

import (
	"context"
	"tts-poc-service/pkg/common/constant"

	"github.com/sirupsen/logrus"
)

var metricLog *logrus.Logger

type Metric struct {
	customMetric map[string]any
}

func (m *Metric) SetKeyValue(key string, value any) {
	m.customMetric[key] = value
}

func (m *Metric) GetMetric() map[string]any {
	return m.customMetric
}

func (m *Metric) GetHashcode() string {
	return m.customMetric[string(constant.CTX_HASHCODE)].(string)
}

type MetricContext struct {
	context.Context
	context.CancelFunc
	*Metric
}

func NewMetricContext() *Metric {
	return &Metric{customMetric: make(map[string]any)}
}

func NewMetricLog(log *logrus.Logger) {
	metricLog = log
}

func SetMetrics(ctx context.Context, data map[string]any) {
	mCtx := TransformCustomMetric(ctx)
	for k, v := range data {
		mCtx.Metric.SetKeyValue(k, v)
	}
}

func SetMetric(ctx context.Context, key string, value any) {
	TransformCustomMetric(ctx).Metric.SetKeyValue(key, value)
}

func CopyContext(ctx context.Context) context.Context {
	mctx := TransformCustomMetric(ctx)
	c := &MetricContext{}
	c.Metric = NewMetricContext()
	c.Context, c.CancelFunc = context.WithCancel(context.Background())
	SetMetrics(c, mctx.GetMetric())
	return c
}

func GetValueContext(ctx context.Context, key string) any {
	return TransformCustomMetric(ctx).customMetric[key]
}

func CancelCustomContext(ctx context.Context) {
	TransformCustomMetric(ctx).CancelFunc()
}

func TransformCustomMetric(ctx context.Context) *MetricContext {
	return ctx.(*MetricContext)
}

func SendMetric(data map[string]any) {
	metricLog.WithFields(data).Debug()
}
