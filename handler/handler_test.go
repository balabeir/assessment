package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/balabeir/assessment/store"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *sql.DB {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return db
}

func TestCreateExpenseHandler(t *testing.T) {
	db := setup(t)
	defer db.Close()
	handler := NewServer(db)
	srv := httptest.NewServer(handler)

	want := store.Expense{
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}
	reqBody, err := json.Marshal(want)
	if err != nil {
		panic(err)
	}
	resp, _ := http.Post(srv.URL+"/expense/", echo.MIMEApplicationJSON, bytes.NewBuffer(reqBody))

	var got store.Expense
	json.NewDecoder(resp.Body).Decode(&got)

	assert := assert.New(t)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	assert.NotEqual(0, got.ID)
	assert.Equal(want.Title, got.Title)
	assert.Equal(want.Amount, got.Amount)
	assert.Equal(want.Note, got.Note)
	assert.Equal(want.Tags, got.Tags)
}

func TestGetExpenseHandler(t *testing.T) {
	db := setup(t)
	defer db.Close()
	handler := NewServer(db)
	srv := httptest.NewServer(handler)

	resp, _ := http.Get(srv.URL + "/expense/1")

	var got store.Expense
	json.NewDecoder(resp.Body).Decode(&got)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, got.ID)
}
