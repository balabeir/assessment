package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/balabeir/assessment/store"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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

	handler.POST("/expense", handler.createExpenseHandler)

	return handler
}

func (h *Handler) createExpenseHandler(c echo.Context) error {
	expense := store.Expense{}
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error bad request",
		})
	}

	err = expense.Create(h.db)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, expense)
}
