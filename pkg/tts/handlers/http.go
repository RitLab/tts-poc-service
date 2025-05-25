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

func (t ttsServer) JoinMP3Files(c echo.Context) (err error) {
	ctx := c.Request().Context()
	// Parse the multipart form, with a maximum memory of 10MB.
	err = c.Request().ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	files := form.File["files"]
	out, err := t.apps.Queries.JoinMp3FilesHandler.Handle(ctx, query.JoinMp3FilesQuery{Files: files})
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t ttsServer) AudioTranscript(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.AudioTranscriptQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	request.File = file
	out, err := t.apps.Queries.AudioTranscriptHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t ttsServer) AudioSummarize(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.AudioSummarizeQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	request.File = file
	out, err := t.apps.Queries.AudioSummarizeHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}
