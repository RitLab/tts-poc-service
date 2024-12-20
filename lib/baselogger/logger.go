package baselogger

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"tts-poc-service/pkg/common/constant"
	pkgMetric "tts-poc-service/pkg/common/metric"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	loggr := &logrus.Logger{
		Out:          io.MultiWriter(os.Stdout),
		ReportCaller: true,
		Level:        logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				s := strings.Split(frame.Function, ".")
				pathFile := s[0]
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				return pathFile, fileName
			},
			TimestampFormat: "2006-01-02 15:04:05.000",
		},
		Hooks: make(map[logrus.Level][]logrus.Hook),
	}

	return &Logger{
		Logger: loggr,
	}
}

func NewMetricLogger(path string) *logrus.Logger {
	return &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename: path,
			MaxAge:   1, //days
		},
		Level: logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			DisableTimestamp: true,
		},
	}
}

func (l *Logger) Hashcode(ctx context.Context) *logrus.Entry {
	return l.Logger.WithContext(ctx).WithField("hashcode", pkgMetric.GetValueContext(ctx, string(constant.CTX_HASHCODE)))
}
