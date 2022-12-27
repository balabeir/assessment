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
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func NewServer(db *sql.DB) *echo.Echo {
	handler := New(db)
	e := echo.New()
	e.POST("/expense", handler.createExpenseHandler)

	return e
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
