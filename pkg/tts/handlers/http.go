package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/common/response"
	pkgUtil "tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/tts/app"
	"tts-poc-service/pkg/tts/app/query"
)

type ttsServer struct {
	apps app.TtService
}

func NewTtsServer(apps app.TtService) ServerInterface {
	return &ttsServer{apps: apps}
}

func (t ttsServer) TextToSpeech(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.GetTextToSpeechFileQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	out, err := t.apps.Queries.GetTextToSpeechHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t ttsServer) ReadTextToSpeech(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.ReadTextToSpeechFileQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	_, err = t.apps.Queries.ReadTextToSpeechHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, nil)
}
