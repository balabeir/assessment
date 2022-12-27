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
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}
	db := database.NewDB(conn)
	database.InitialDB(db)

	handler := database.NewServer(db.DB)

	// start server
	go func() {
		if err := handler.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			handler.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown(handler)
}

func shutdown(h *database.Handler) {
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
