package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/balabeir/assessment/database"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

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
	e.GET("/expense/:id", handler.getExpenseHandler)
	e.PUT("/expense/:id", handler.updateExpenseHandler)

	return e
}

func (h *Handler) createExpenseHandler(c echo.Context) error {
	expense := database.Expense{}
	err := c.Bind(&expense)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})
	}

	err = expense.Create(h.db)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "internal server error"})
	}

	return c.JSON(http.StatusCreated, expense)
}

func (h *Handler) getExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})

	}

	expense := database.Expense{ID: id}
	err = expense.Get(h.db)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})

	} else if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "internal server error"})
	}

	return c.JSON(http.StatusOK, expense)
}

func (h *Handler) updateExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})
	}

	expense := database.Expense{}
	err = c.Bind(&expense)
	if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})
	}

	expense.ID = id

	err = expense.Update(h.db)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, Err{Message: "error bad request"})

	} else if err != nil {
		c.Echo().Logger.Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "internal server error"})
	}

	return c.JSON(http.StatusOK, expense)
}
