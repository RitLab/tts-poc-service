package pkg_error

import (
	"embed"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

//go:embed error_properties.json
var error_properties embed.FS
var errorData map[string]ErrorData

type ErrorData struct {
	Status int `json:"status,omitempty"`
	Body   any `json:"body,omitempty"`
}

func constructErrorData() {
	filename := "error_properties.json"
	errorProperties, _ := error_properties.ReadFile(filename)
	json.Unmarshal(errorProperties, &errorData)
}

func CreateError(c echo.Context, errorKey string) error {
	if errorData == nil {
		constructErrorData()
	}

	if value, found := errorData[errorKey]; found {
		return c.JSON(value.Status, value.Body)
	}
	return c.JSON(errorData[GENERAL_ERROR].Status, errorData[GENERAL_ERROR].Body)
}

func CreateCustomError(c echo.Context, status int, slug, msg string) error {
	body := map[string]string{
		"slug":    slug,
		"message": msg,
	}
	return c.JSON(status, body)
}
