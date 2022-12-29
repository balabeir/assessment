package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balabeir/assessment/store"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreateExpenseHandler(t *testing.T) {
	expense := store.Expense{
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

	db, mock := setupDB(t)
	defer db.Close()
	handler := New(db)

	mock.ExpectExec("INSERT INTO expenses").
		WithArgs(expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ex := expense
	ex.ID = 1
	expected, _ := json.Marshal(ex)

	err := handler.createExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, string(expected), strings.TrimSpace(res.Body.String()))
	}
}

func TestGetExpenseHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/expense/1", strings.NewReader(""))
	res := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, res)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	db, mock := setupDB(t)
	defer db.Close()
	handler := New(db)

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "Bob", 20, "testing", pq.Array([]string{"foo", "bar"}))

	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(mockRows)

	expense := store.Expense{
		ID:     1,
		Title:  "Bob",
		Amount: 20,
		Note:   "testing",
		Tags:   []string{"foo", "bar"},
	}
	expected, _ := json.Marshal(expense)

	err := handler.getExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(expected), strings.TrimSpace(res.Body.String()))
	}

	// db := setup(t)
	// defer db.Close()
	// handler := NewServer(db)
	// srv := httptest.NewServer(handler)

	// resp, _ := http.Get(srv.URL + "/expense/1")

	// var got store.Expense
	// json.NewDecoder(resp.Body).Decode(&got)

	// assert := assert.New(t)
	// assert.Equal(http.StatusOK, resp.StatusCode)
	// assert.Equal(1, got.ID)
}
