package response

import (
	"github.com/labstack/echo/v4"
	pkgUtil "tts-poc-service/pkg/common/utils"
)

type response struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
}

func SuccessResponse(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, response{
		Status:  status,
		Message: "success",
		Data:    data,
	})
}

type HealthCheckResponse struct {
	StartTime string `json:"start_time"`
	Status    string `json:"status"`
}

type Pagination struct {
	PerPage      int  `json:"per_page"`
	TotalRecords int  `json:"total_records"`
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	NextPage     *int `json:"next_page"`
	PrevPage     *int `json:"prev_page"`
}

func NewPagination(page, perPage, totalRecords int) *Pagination {
	var nextPage, prevPage, totalPages int
	totalPages = totalRecords / perPage
	if totalRecords%perPage != 0 {
		totalPages++
	}

	if page > 1 {
		prevPage = page - 1
	}

	if perPage*page < totalRecords {
		nextPage = page + 1
	}
	return &Pagination{
		PerPage:      perPage,
		TotalRecords: totalRecords,
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     pkgUtil.IfThenElse(nextPage != 0, &nextPage, nil),
		PrevPage:     pkgUtil.IfThenElse(prevPage != 0, &prevPage, nil),
	}
}
