package handler

import (
	"github.com/labstack/echo/v4"
	"time"
	"tts-poc-service/pkg/common/response"
)

type HealthCheckHandler interface {
	HealthCheck(c echo.Context) error
}

type healthCheckHandler struct {
	StartTime time.Time
}

func NewHealthCheckHandler(start time.Time) HealthCheckHandler {
	return &healthCheckHandler{StartTime: start}
}

func (h *healthCheckHandler) HealthCheck(c echo.Context) error {
	resp := &response.HealthCheckResponse{
		StartTime: h.StartTime.Format("2006-01-02 15:04:05 Monday"),
		Status:    "OK",
	}
	return c.JSON(200, resp)
}
