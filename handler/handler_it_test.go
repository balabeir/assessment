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

	"github.com/balabeir/assessment/store"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func TestITCreateExpense(t *testing.T) {
	db, _ := sql.Open("postgres", "postgres://xdkhrnfq:ri_5P5A5v_Z-uGAoeyaLset9oWhN24xv@babar.db.elephantsql.com/xdkhrnfq")

	expense := store.Expense{
		ID:     1,
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}
	reqBody, _ := json.Marshal(expense)

	handler := New(db)
	e := echo.New()
	e.POST("/expense", handler.createExpenseHandler)

	srv := httptest.NewServer(e)
	resp, _ := http.Post(srv.URL+"/expense", echo.MIMEApplicationJSON, bytes.NewBuffer(reqBody))

	var got store.Expense
	json.NewDecoder(resp.Body).Decode(&got)

	if assert.Equal(t, http.StatusCreated, resp.StatusCode) {
		assert.NotEqual(t, 0, got.ID)
		assert.Equal(t, expense.Title, got.Title)
		assert.Equal(t, expense.Amount, got.Amount)
		assert.Equal(t, expense.Note, got.Note)
		assert.Equal(t, expense.Tags, got.Tags)
	}
}
