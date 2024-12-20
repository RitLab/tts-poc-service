package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/common/response"
	"tts-poc-service/pkg/config/app"
)

type ConfigServerHandler interface {
	GetConfig(c echo.Context) error
	ReloadConfig(c echo.Context) error
}

type configHttpHandler struct {
	logger *baselogger.Logger
	app    app.ConfigService
}

func NewConfigHttpHandler(logger *baselogger.Logger, app app.ConfigService) ConfigServerHandler {
	return &configHttpHandler{logger: logger, app: app}
}

func (cf *configHttpHandler) GetConfig(c echo.Context) (err error) {
	cfg, _ := cf.app.Queries.GetConfigHandler.Handle(c.Request().Context(), config.Config)
	return c.JSON(http.StatusOK, cfg)
}

func (cf *configHttpHandler) ReloadConfig(c echo.Context) (err error) {
	cf.app.Commands.ReloadConfigHandler.Handle(c.Request().Context(), config.Config)
	return response.SuccessResponse(c, http.StatusOK, nil)
}
