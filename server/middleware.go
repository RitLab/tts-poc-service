package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/common/constant"
	pkgMetric "tts-poc-service/pkg/common/metric"
	pkgUtil "tts-poc-service/pkg/common/utils"
)

const (
	__header_auth     = "Authorization"
	__header_key_room = "x-key-room"
	__header_hashcode = "X-Hashcode"
)

type IDTokenHandler struct {
	JwkURL     string
	EmailKey   string
	UserKey    string
	CountryKey string
}

type MiddlewareHTTP interface {
	Middleware(next http.Handler) http.Handler
	MiddlewareWithLogger() middleware.RequestLoggerConfig
}

type _http struct {
	*baselogger.Logger
}

func New(log *baselogger.Logger) MiddlewareHTTP {
	return &_http{Logger: log}
}

func (h _http) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &pkgMetric.MetricContext{
			Context: r.Context(),
			Metric:  pkgMetric.NewMetricContext(),
		}

		hx := r.Header.Get(__header_hashcode)
		hc := pkgUtil.IfThenElse(len(hx) == 0, generateHashcode(), hx)
		ctx.Metric.SetKeyValue(string(constant.CTX_HASHCODE), hc)

		r = r.WithContext(ctx)
		w.Header().Set(__header_hashcode, hc)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func (h _http) MiddlewareWithLogger() middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogHost:      true,
		LogMethod:    true,
		LogUserAgent: true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogReferer:   true,
		LogRequestID: true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()
			if strings.Compare(v.URI, "/api/tts/health-check") == 0 {
				return nil
			}
			if strings.Contains(v.URI, "/api/tts/swagger") {
				return nil
			}
			if v.Error == nil {
				h.Hashcode(ctx).Logger.WithContext(ctx).WithFields(logrus.Fields{
					"uri":     v.URI,
					"status":  v.Status,
					"latency": v.Latency,
				}).Info("request")
			} else {
				h.Hashcode(ctx).Logger.WithContext(ctx).WithFields(logrus.Fields{
					"uri":     v.URI,
					"status":  v.Status,
					"latency": v.Latency,
					"error":   v.Error,
				}).Error("request error")
			}
			pkgMetric.SendMetric(map[string]any{
				"uri":        v.URI,
				"status":     v.Status,
				"host":       v.Host,
				"method":     v.Method,
				"user_agent": v.UserAgent,
				"latency":    v.Latency,
				"remote_ip":  v.RemoteIP,
				"referer":    v.Referer,
				"duration":   time.Since(v.StartTime),
			})
			return nil
		},
	}
}

func generateHashcode() string {
	return "(" + strings.Split(uuid.New().String(), "-")[0] + ")"
}
