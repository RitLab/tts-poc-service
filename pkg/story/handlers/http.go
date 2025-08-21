package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/common/response"
	pkgUtil "tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/story/app"
	"tts-poc-service/pkg/story/app/query"
)

type storyServer struct {
	apps app.StoryService
}

func NewStoryServer(apps app.StoryService) ServerInterface {
	return &storyServer{apps: apps}
}

func (t storyServer) GenerateStory(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.GenerateStoryBookQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	out, err := t.apps.Queries.GenerateStoryBookHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t storyServer) GetStory(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.GetStoryBookQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	out, err := t.apps.Queries.GetStoryBookHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t storyServer) ListStory(c echo.Context) (err error) {
	ctx := c.Request().Context()

	out, err := t.apps.Queries.ListStoryBookHandler.Handle(ctx, query.ListStoryBookQuery{})
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}
