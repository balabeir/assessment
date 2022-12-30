package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/balabeir/assessment/database"
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

	g := e.Group("/expense")
	g.GET("/:id", handler.getExpenseHandler)
	g.POST("/", handler.createExpenseHandler)

	return e
}

func (h *Handler) createExpenseHandler(c echo.Context) error {
	expense := database.Expense{}
	err := c.Bind(&expense)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error bad request",
		})
	}

	err = expense.Create(h.db)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, expense)
}

func (h *Handler) getExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "error bad request",
		})
	}

	expense := database.Expense{ID: id}
	err = expense.Get(h.db)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "bad request",
		})
	} else if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, expense)
}
