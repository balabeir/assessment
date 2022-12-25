package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/balabeir/assessment/handler"
	"github.com/labstack/echo"
)

func main() {
	h := handler.New()

	// start server
	go func() {
		if err := h.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			h.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown(h)
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
