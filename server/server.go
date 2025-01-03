package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/lib/storage"
	"tts-poc-service/lib/validator"
	pkgMetric "tts-poc-service/pkg/common/metric"
	configApp "tts-poc-service/pkg/config/app"
	configHandler "tts-poc-service/pkg/config/handler"
	healthHandler "tts-poc-service/pkg/health_check/handler"
	"tts-poc-service/pkg/tts/app"
	"tts-poc-service/pkg/tts/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	DBInitFile = "./migration/init.sql"
)

type Handler struct {
	HealthCheck  healthHandler.HealthCheckHandler
	ConfigServer configHandler.ConfigServerHandler
	TtsService   handlers.ServerInterface
}

func setHandler(dep Dependency) Handler {
	return Handler{
		HealthCheck:  healthHandler.NewHealthCheckHandler(time.Now()),
		ConfigServer: configHandler.NewConfigHttpHandler(dep.logger, configApp.NewConfigService(dep.logger)),
		TtsService:   handlers.NewTtsServer(app.NewTtsService(dep.logger, dep.player, dep.storage)),
	}
}

type Dependency struct {
	logger  *baselogger.Logger
	player  htgo.Player
	storage storage.Storage
	val     *validator.Validator
}

func newDependency(ctx context.Context) Dependency {
	logger := baselogger.NewLogger()
	metricLog := baselogger.NewMetricLogger("/tmp/be/metric.log")
	pkgMetric.NewMetricLog(metricLog)

	val := validator.New()
	config.InitConfig(ctx, logger)

	s3 := storage.NewMinioHandler(logger)
	player := htgo.Player{}

	return Dependency{
		logger:  logger,
		player:  player,
		storage: s3,
		val:     val,
	}
}

type Server interface {
	Start() error
	HandleShutdown(ctx context.Context) context.Context
}

type server struct {
	Dependency
	*http.Server
}

func NewServer(ctx context.Context) Server {
	dep := newDependency(ctx)
	srvc := echo.New()
	hndler := setHandler(dep)

	mid := New(dep.logger)
	setMiddleware(srvc, dep, mid)

	root := srvc.Group("/api/tts")

	// Serve OpenAPI Specification
	root.GET("/openapi.yaml", func(c echo.Context) error {
		return c.File("api/openapi/tts.yaml")
	})

	// Serve Swagger UI Static Files
	root.Static("/swagger-ui", "public/swagger-ui")

	// Redirect to Swagger UI with OpenAPI Specification URL
	root.GET("/docs", func(c echo.Context) error {
		redirectURL := "/api/tts/swagger-ui/index.html?url=/api/tts/openapi.yaml"
		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	})

	// general api
	root.GET("/health-check", hndler.HealthCheck.HealthCheck)
	root.POST("/config/reload", hndler.ConfigServer.ReloadConfig)
	root.GET("/config", hndler.ConfigServer.GetConfig)

	// tts
	root.POST("", hndler.TtsService.TextToSpeech)
	root.POST("/read", hndler.TtsService.ReadTextToSpeech)

	srvr := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		ReadTimeout:  time.Duration(config.Config.Server.ReadTimeout),
		Handler:      srvc,
		WriteTimeout: time.Duration(config.Config.Server.ReadTimeout),
	}

	return &server{dep, srvr}
}

func (s *server) Start() error {
	s.logger.Infof(fmt.Sprintf("Service started at port %d...", config.Config.Server.Port))
	return s.ListenAndServe()
}

func (s *server) HandleShutdown(ctx context.Context) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		if err := s.Shutdown(ctx); err != nil {
			s.logger.Errorf("could not gracefully shutdown the api server")
		} else {
			s.logger.Info("api server is shutting down")
		}
	}()
	return ctx
}

func setMiddleware(router *echo.Echo, dep Dependency, mid MiddlewareHTTP) {
	router.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))

	router.Validator = dep.val
	router.Use(middleware.CORS())
	router.Use(echo.WrapMiddleware(mid.Middleware))
	router.Use(middleware.RequestLoggerWithConfig(mid.MiddlewareWithLogger()))
	router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			dep.logger.Logger.Error("Panic: ", err, " stack: ", string(stack))
			return nil
		},
	}))
	router.Use(contentSecurityPolicy)
}

func contentSecurityPolicy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Security-Policy", "style-src 'self' 'unsafe-inline';")
		return next(c)
	}
}
