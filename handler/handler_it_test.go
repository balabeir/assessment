//go:build integration
// +build integration

package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

const serverPort = 80

func TestITCreateExpense(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		handler := New(db)
		e.GET("/expense/:id", handler.getExpenseHandler)

		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}
