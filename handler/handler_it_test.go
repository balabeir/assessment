//go:build integration
// +build integration

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	handler := New(db)

	err := handler.createExpenseHandler(c)
	var got database.Expense
	json.NewDecoder(res.Body).Decode(&got)

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

	req := httptest.NewRequest(http.MethodGet, "/expenses/1", strings.NewReader(""))
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	c.SetPath("/expenses/:id")
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
	var got database.Expense
	json.NewDecoder(res.Body).Decode(&got)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expected.ID, got.ID)
		assert.Equal(t, expected.Title, got.Title)
		assert.Equal(t, expected.Amount, got.Amount)
		assert.Equal(t, expected.Note, got.Note)
		assert.Equal(t, expected.Tags, got.Tags)
	}
}

func TestITUpdateExpense(t *testing.T) {
	db := setupDBIntegration(t)
	defer db.Close()

	expected := database.Expense{
		Title:  "Malee",
		Amount: 100,
		Note:   "just testing",
		Tags:   []string{"banana", "orange"},
	}
	reqBody, _ := json.Marshal(expected)

	req := httptest.NewRequest(http.MethodPut, "/expenses/2", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")
	handler := New(db)
	id, _ := strconv.Atoi(c.Param("id"))
	expected.ID = id

	err := handler.updateExpenseHandler(c)
	var got database.Expense
	json.NewDecoder(res.Body).Decode(&got)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expected.ID, got.ID)
		assert.Equal(t, expected.Title, got.Title)
		assert.Equal(t, expected.Amount, got.Amount)
		assert.Equal(t, expected.Note, got.Note)
		assert.Equal(t, expected.Tags, got.Tags)
	}
}

func TestITGetExpenseLists(t *testing.T) {
	db := setupDBIntegration(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	handler := New(db)

	err := handler.getExpenseListsHandler(c)
	var got []database.Expense
	json.NewDecoder(res.Body).Decode(&got)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Less(t, 0, len(got))
	}
}
