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

	// e := echo.New()
	// req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)
	// h := New(db)

	// expected, _ := json.Marshal(expense)

	// Assertions
	// if assert.NoError(t, h.createExpenseHandler(c)) {
	// assert.Equal(t, http.StatusCreated, rec.Code)
	// assert.Equal(t, string(expected), rec.Body.String())
	// }
	// eh := echo.New()
	// go func(e *echo.Echo) {
	// 	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	handler := New(db)
	// 	e.GET("/expense/:id", handler.getExpenseHandler)

	// 	e.Start(fmt.Sprintf(":%d", serverPort))
	// }(eh)
	// for {
	// 	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	if conn != nil {
	// 		conn.Close()
	// 		break
	// 	}
	// }
}
