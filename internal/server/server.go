// server  - describes a container for pointer http.Server
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ekvo/golang-chi-postgres-api/internal/config"
)

const TimeShutServer = 10 * time.Second

// Connect - wrapper http.Server
type Connect struct {
	*http.Server
}

func NewServer(srv *http.Server) *Connect {
	return &Connect{Server: srv}
}

func Init(cfg *config.Config, r http.Handler) *Connect {
	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.ServerHost),
		Handler: r,
	}
	return NewServer(srv)
}

func (c *Connect) ListenAndServeAndShut(ctx context.Context, timeShut time.Duration) error {
	go func() {
		log.Print("server: Listen and serve - start\n")
		if err := c.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: HTTP server error - %v", err)
		}
		log.Print("server: stopped serving\n")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, timeShut)
	defer shutdownRelease()

	if err := c.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server: HTTP shutdown error - %w", err)
	}
	log.Print("server: graceful shutdown complete\n")
	return nil
}
