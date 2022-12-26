package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// type Handler struct {
// 	DB *sql.DB
// }

func New() *echo.Echo {
	e := echo.New()
	// write handler here
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "OK",
		})
	})

	return e
}
