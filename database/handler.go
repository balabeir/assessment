package database

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	*echo.Echo
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{echo.New(), db}
}

func NewServer(db *sql.DB) *Handler {
	handler := NewHandler(db)

	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	handler.POST("/expense", handler.createExpenseHandler)

	return handler
}

func (h *Handler) createExpenseHandler(c echo.Context) error {
	expense := Expense{}
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error bad request",
		})
	}

	err = expense.Create(h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, expense)
}
