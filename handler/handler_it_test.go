//go:build integration
// +build integration

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/balabeir/assessment/database"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDBIntegration(t *testing.T) *sql.DB {
	t.Parallel()
	db, err := sql.Open("postgres", "postgresql://test:test@db/it-db?sslmode=disable")
	assert.NoError(t, err)

	return db
}

func TestITCreateExpense(t *testing.T) {
	db := setupDBIntegration(t)
	defer db.Close()

	expense := database.Expense{
		ID:     2,
		Title:  "John",
		Amount: 100,
		Note:   "paid",
		Tags:   []string{"pet", "market"},
	}
	reqBody, _ := json.Marshal(expense)

	req := httptest.NewRequest(http.MethodPost, "/expense", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	handler := New(db)

	err := handler.createExpenseHandler(c)
	assert.NoError(t, err)
	var got database.Expense
	err = json.NewDecoder(res.Body).Decode(&got)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEqual(t, 0, got.ID)
		assert.Equal(t, expense.Title, got.Title)
		assert.Equal(t, expense.Amount, got.Amount)
		assert.Equal(t, expense.Note, got.Note)
		assert.Equal(t, expense.Tags, got.Tags)
	}
}

func TestITGetExpense(t *testing.T) {
	db := setupDBIntegration(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/expense/1", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	handler := New(db)

	expected := database.Expense{
		ID:     1,
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}

	err := handler.getExpenseHandler(c)
	assert.NoError(t, err)
	var got database.Expense
	err = json.NewDecoder(res.Body).Decode(&got)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expected.ID, got.ID)
		assert.Equal(t, expected.Title, got.Title)
		assert.Equal(t, expected.Amount, got.Amount)
		assert.Equal(t, expected.Note, got.Note)
		assert.Equal(t, expected.Tags, got.Tags)
	}
}
