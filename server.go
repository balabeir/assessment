package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/balabeir/assessment/database"
	"github.com/balabeir/assessment/handler"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	database.Initial(conn)
	e := handler.NewServer(conn)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// start server
	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	defer cancel()
}
