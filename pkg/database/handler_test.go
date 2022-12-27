package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseHandler(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	handler := NewHandler(db)
	srv := httptest.NewServer(handler.E)

	want := Expense{
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}
	reqBody, err := json.Marshal(want)
	if err != nil {
		panic(err)
	}
	resp, _ := http.Post(srv.URL+"/expense", echo.MIMEApplicationJSON, bytes.NewBuffer(reqBody))

	var got Expense
	json.NewDecoder(resp.Body).Decode(&got)

	assert := assert.New(t)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	assert.NotEqual(0, got.ID)
	assert.Equal(want.Title, got.Title)
	assert.Equal(want.Amount, got.Amount)
	assert.Equal(want.Note, got.Note)
	assert.Equal(want.Tags, got.Tags)
}
