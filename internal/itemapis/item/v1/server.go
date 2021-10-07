package v1

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func RunServer(e *echo.Echo, srv Service, addr string) error {

	errChan := make(chan error)

	// Set log level
	e.Logger.SetLevel(log.INFO)

	// Get routes
	initRoutes(e, srv)

	// Start server
	go func() {
		e.Logger.Infof("Item service is running")
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		e.Logger.Infof("got %v signal, graceful shutdown server", <-c)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	return <-errChan
}
