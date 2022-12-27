package database

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	E  *echo.Echo
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	e := echo.New()
	handler := &Handler{
		E:  e,
		db: db,
	}

	handler.E.Use(middleware.Logger())
	handler.E.Use(middleware.Recover())

	handler.E.POST("/expense", handler.createExpenseHandler)

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

	err = expense.CreateExpense(h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, expense)
}
