package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/common/response"
	pkgUtil "tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/pdf/app"
	"tts-poc-service/pkg/pdf/app/command"
	"tts-poc-service/pkg/pdf/app/query"
)

type pdfServer struct {
	apps app.PdfService
}

func NewPdfServer(apps app.PdfService) ServerInterface {
	return &pdfServer{apps: apps}
}

func (t pdfServer) JoinPdfFiles(c echo.Context) (err error) {
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
	out, err := t.apps.Queries.JoinPdfFilesHandler.Handle(ctx, query.JoinPdfFilesQuery{Files: files})
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t pdfServer) SignPdfFile(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request query.SignPdfFileQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	request.File = file
	out, err := t.apps.Queries.SignPdfFileHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, out)
}

func (t pdfServer) VerifyPdfFile(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request command.VerifyPdfFileQuery
	if err = pkgUtil.BindRequestAndValidate(c, &request); err != nil {
		return pkgError.CreateCustomError(c, http.StatusBadRequest, "bad-request", err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	request.File = file
	err = t.apps.Commands.VerifyPdfFileHandler.Handle(ctx, request)
	if err != nil {
		return pkgError.CreateError(c, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, nil)
}
