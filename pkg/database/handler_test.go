package database

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	InitialDB()
	e := NewHandler()
	srv := httptest.NewServer(e)

	want := Expense{
		Title:  "Bob",
		Amount: 10,
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
