package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/common/response"
	pkgUtil "tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/support/app"
	"tts-poc-service/pkg/support/app/command"
)

type supportServer struct {
	apps app.SupportService
}

func NewSupportServer(apps app.SupportService) ServerInterface {
	return &supportServer{apps: apps}
}

func (t supportServer) InsertSupport(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request command.InsertSupportCommand
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	err = t.apps.Commands.InsertSupportHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, nil)
}
