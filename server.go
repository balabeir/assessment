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

	"github.com/balabeir/assessment/handler"
	"github.com/balabeir/assessment/store"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	store.InitialDB(conn)
	handler := handler.NewServer(conn)
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	// start server
	go func() {
		if err := handler.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			handler.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown(handler)
}

func shutdown(h *handler.Handler) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil {
		h.Logger.Fatal(err)
	}
}
