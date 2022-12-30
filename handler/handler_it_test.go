//go:build integration
// +build integration

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/balabeir/assessment/database"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestITCreateExpense(t *testing.T) {
	db, err := sql.Open("postgres", "postgresql://test:test@db/it-db?sslmode=disable")
	assert.NoError(t, err)

	expense := database.Expense{
		ID:     1,
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}
	reqBody, _ := json.Marshal(expense)

	req := httptest.NewRequest(http.MethodPost, "/expense/", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	handler := New(db)

	err = handler.createExpenseHandler(c)
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
