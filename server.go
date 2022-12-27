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

	"github.com/balabeir/assessment/pkg/database"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}
	database.InitialDB(db)
	handler := database.NewHandler(db)

	// start server
	go func() {
		if err := handler.E.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			handler.E.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown(handler.E)
}

func shutdown(h *echo.Echo) {
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
